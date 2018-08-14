package price_feeds

import (
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/price_feeds"
)

type MockPriceFeedRepository struct {
	createErr                 error
	missingHeaders            []core.Header
	missingHeadersErr         error
	passedEndingBlockNumber   int64
	passedModel               price_feeds.PriceFeedModel
	passedStartingBlockNumber int64
}

func (repository *MockPriceFeedRepository) SetCreateErr(err error) {
	repository.createErr = err
}

func (repository *MockPriceFeedRepository) SetMissingHeadersErr(err error) {
	repository.missingHeadersErr = err
}

func (repository *MockPriceFeedRepository) SetMissingHeaders(headers []core.Header) {
	repository.missingHeaders = headers
}

func (repository *MockPriceFeedRepository) Create(model price_feeds.PriceFeedModel) error {
	repository.passedModel = model
	return repository.createErr
}

func (repository *MockPriceFeedRepository) MissingHeaders(startingBlockNumber, endingBlockNumber int64) ([]core.Header, error) {
	repository.passedStartingBlockNumber = startingBlockNumber
	repository.passedEndingBlockNumber = endingBlockNumber
	return repository.missingHeaders, repository.missingHeadersErr
}

func (repository *MockPriceFeedRepository) AssertCreateCalledWith(model price_feeds.PriceFeedModel) {
	Expect(repository.passedModel).To(Equal(model))
}

func (repository *MockPriceFeedRepository) AssertMissingHeadersCalledwith(startingBlockNumber, endingBlockNumber int64) {
	Expect(repository.passedStartingBlockNumber).To(Equal(startingBlockNumber))
	Expect(repository.passedEndingBlockNumber).To(Equal(endingBlockNumber))
}
