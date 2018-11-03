// Copyright 2018 Vulcanize
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package converter_test

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/examples/constants"
	"github.com/vulcanize/vulcanizedb/examples/generic/helpers"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/omni/converter"
	"github.com/vulcanize/vulcanizedb/pkg/omni/parser"
	"github.com/vulcanize/vulcanizedb/pkg/omni/types"
)

var mockEvent = core.WatchedEvent{
	LogID:       1,
	Name:        constants.TransferEvent.String(),
	BlockNumber: 5488076,
	Address:     constants.TusdContractAddress,
	TxHash:      "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad6546ae",
	Index:       110,
	Topic0:      constants.TransferEvent.Signature(),
	Topic1:      "0x000000000000000000000000000000000000000000000000000000000000af21",
	Topic2:      "0x9dd48110dcc444fdc242510c09bbbbe21a5975cac061d82f7b843bce061ba391",
	Topic3:      "",
	Data:        "0x000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc200000000000000000000000089d24a6b4ccb1b6faa2625fe562bdd9a23260359000000000000000000000000000000000000000000000000392d2e2bda9c00000000000000000000000000000000000000000000000000927f41fa0a4a418000000000000000000000000000000000000000000000000000000000005adcfebe",
}

var _ = Describe("Converter Test", func() {

	It("Converts watched event log to mapping of event input names to values", func() {
		p := parser.NewParser("")
		err := p.Parse(constants.TusdContractAddress)
		Expect(err).ToNot(HaveOccurred())

		info := types.ContractInfo{
			Name:          "TrueUSD",
			Address:       constants.TusdContractAddress,
			Abi:           p.Abi(),
			ParsedAbi:     p.ParsedAbi(),
			StartingBlock: 5197514,
			Events:        p.GetEvents(),
			Methods:       p.GetMethods(),
		}

		event := info.Events["Transfer"]

		info.GenerateFilters([]string{"Transfer"})
		c := converter.NewConverter(info)
		err = c.Convert(mockEvent, event)
		Expect(err).ToNot(HaveOccurred())

		from := common.HexToAddress("0x000000000000000000000000000000000000000000000000000000000000af21")
		to := common.HexToAddress("0x9dd48110dcc444fdc242510c09bbbbe21a5975cac061d82f7b843bce061ba391")
		value := helpers.BigFromString("1097077688018008265106216665536940668749033598146")

		v := event.Logs[1].Values["value"].(*big.Int)

		Expect(event.Logs[1].Values["to"].(common.Address)).To(Equal(to))
		Expect(event.Logs[1].Values["from"].(common.Address)).To(Equal(from))
		Expect(v.String()).To(Equal(value.String()))
	})
})