// VulcanizeDB
// Copyright © 2018 Vulcanize

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.

// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package geth_test

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	vulcCore "github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/geth"
	"github.com/vulcanize/vulcanizedb/pkg/geth/converters/cold_db"
)

var _ = Describe("Geth blockchain", func() {
	var mockClient *fakes.MockEthClient
	var mockRpcClient *fakes.MockRpcClient
	var node vulcCore.Node
	var blockChain *geth.BlockChain

	BeforeEach(func() {
		mockClient = fakes.NewMockEthClient()
		mockRpcClient = fakes.NewMockRpcClient()
		node = vulcCore.Node{}
		blockChain = geth.NewBlockChain(mockClient, mockRpcClient, node, cold_db.NewColdDbTransactionConverter())
	})

	Describe("getting a block", func() {
		It("fetches block from ethClient", func() {
			mockClient.SetBlockByNumberReturnBlock(types.NewBlockWithHeader(&types.Header{}))
			blockNumber := int64(100)

			_, err := blockChain.GetBlockByNumber(blockNumber)

			Expect(err).NotTo(HaveOccurred())
			mockClient.AssertBlockByNumberCalledWith(context.Background(), big.NewInt(blockNumber))
		})

		It("returns err if ethClient returns err", func() {
			mockClient.SetBlockByNumberErr(fakes.FakeError)

			_, err := blockChain.GetBlockByNumber(100)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(fakes.FakeError))
		})
	})

	Describe("getting a header", func() {
		Describe("default/mainnet", func() {
			It("fetches header from ethClient", func() {
				blockNumber := int64(100)
				mockClient.SetHeaderByNumberReturnHeader(&types.Header{Number: big.NewInt(blockNumber)})

				_, err := blockChain.GetHeaderByNumber(blockNumber)

				Expect(err).NotTo(HaveOccurred())
				mockClient.AssertHeaderByNumberCalledWith(context.Background(), big.NewInt(blockNumber))
			})

			It("returns err if ethClient returns err", func() {
				mockClient.SetHeaderByNumberErr(fakes.FakeError)

				_, err := blockChain.GetHeaderByNumber(100)

				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError(fakes.FakeError))
			})
		})

		Describe("POA/Kovan", func() {
			It("fetches header from rpcClient", func() {
				node.NetworkID = vulcCore.KOVAN_NETWORK_ID
				blockNumber := hexutil.Big(*big.NewInt(123))
				mockRpcClient.SetReturnPOAHeader(vulcCore.POAHeader{Number: &blockNumber})
				blockChain = geth.NewBlockChain(mockClient, mockRpcClient, node, cold_db.NewColdDbTransactionConverter())

				_, err := blockChain.GetHeaderByNumber(100)

				Expect(err).NotTo(HaveOccurred())
				mockRpcClient.AssertCallContextCalledWith(context.Background(), &vulcCore.POAHeader{}, "eth_getBlockByNumber")
			})

			It("returns err if rpcClient returns err", func() {
				node.NetworkID = vulcCore.KOVAN_NETWORK_ID
				mockRpcClient.SetCallContextErr(fakes.FakeError)
				blockChain = geth.NewBlockChain(mockClient, mockRpcClient, node, cold_db.NewColdDbTransactionConverter())

				_, err := blockChain.GetHeaderByNumber(100)

				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError(fakes.FakeError))
			})

			It("returns error if returned header is empty", func() {
				node.NetworkID = vulcCore.KOVAN_NETWORK_ID
				blockChain = geth.NewBlockChain(mockClient, mockRpcClient, node, cold_db.NewColdDbTransactionConverter())

				_, err := blockChain.GetHeaderByNumber(100)

				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError(geth.ErrEmptyHeader))
			})
		})
	})

	Describe("getting logs with default FilterQuery", func() {
		It("fetches logs from ethClient", func() {
			mockClient.SetFilterLogsReturnLogs([]types.Log{{}})
			contract := vulcCore.Contract{Hash: common.BytesToHash([]byte{1, 2, 3, 4, 5}).Hex()}
			startingBlockNumber := big.NewInt(1)
			endingBlockNumber := big.NewInt(2)

			_, err := blockChain.GetLogs(contract, startingBlockNumber, endingBlockNumber)

			Expect(err).NotTo(HaveOccurred())
			expectedQuery := ethereum.FilterQuery{
				FromBlock: startingBlockNumber,
				ToBlock:   endingBlockNumber,
				Addresses: []common.Address{common.HexToAddress(contract.Hash)},
			}
			mockClient.AssertFilterLogsCalledWith(context.Background(), expectedQuery)
		})

		It("returns err if ethClient returns err", func() {
			mockClient.SetFilterLogsErr(fakes.FakeError)
			contract := vulcCore.Contract{Hash: common.BytesToHash([]byte{1, 2, 3, 4, 5}).Hex()}
			startingBlockNumber := big.NewInt(1)
			endingBlockNumber := big.NewInt(2)

			_, err := blockChain.GetLogs(contract, startingBlockNumber, endingBlockNumber)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(fakes.FakeError))
		})
	})

	Describe("getting logs with a custom FilterQuery", func() {
		It("fetches logs from ethClient", func() {
			mockClient.SetFilterLogsReturnLogs([]types.Log{{}})
			address := common.HexToAddress("0x")
			startingBlockNumber := big.NewInt(1)
			endingBlockNumber := big.NewInt(2)
			topic := common.HexToHash("0x")
			query := ethereum.FilterQuery{
				FromBlock: startingBlockNumber,
				ToBlock:   endingBlockNumber,
				Addresses: []common.Address{address},
				Topics:    [][]common.Hash{{topic}},
			}

			_, err := blockChain.GetEthLogsWithCustomQuery(query)

			Expect(err).NotTo(HaveOccurred())
			mockClient.AssertFilterLogsCalledWith(context.Background(), query)
		})

		It("returns err if ethClient returns err", func() {
			mockClient.SetFilterLogsErr(fakes.FakeError)
			contract := vulcCore.Contract{Hash: common.BytesToHash([]byte{1, 2, 3, 4, 5}).Hex()}
			startingBlockNumber := big.NewInt(1)
			endingBlockNumber := big.NewInt(2)
			query := ethereum.FilterQuery{
				FromBlock: startingBlockNumber,
				ToBlock:   endingBlockNumber,
				Addresses: []common.Address{common.HexToAddress(contract.Hash)},
				Topics:    nil,
			}

			_, err := blockChain.GetEthLogsWithCustomQuery(query)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(fakes.FakeError))
		})
	})

	Describe("getting the most recent block number", func() {
		It("fetches latest header from ethClient", func() {
			blockNumber := int64(100)
			mockClient.SetHeaderByNumberReturnHeader(&types.Header{Number: big.NewInt(blockNumber)})

			result := blockChain.LastBlock()

			mockClient.AssertHeaderByNumberCalledWith(context.Background(), nil)
			Expect(result).To(Equal(big.NewInt(blockNumber)))
		})
	})
})
