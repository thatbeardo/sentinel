package mocks

import "errors"

// ResourceWithoutParent returns
func ResourceWithoutParent() *Result {
	mockResult := &Result{}
	mockRecord := &Record{}
	mockRecord.On("GetByIndex", 0).Return("test-resource")
	mockRecord.On("GetByIndex", 1).Return("test-source-id")
	mockRecord.On("GetByIndex", 2).Return("test-id")
	mockRecord.On("GetByIndex", 3).Return(nil)
	mockResult.On("Next").Return(true).Once()
	mockResult.On("Next").Return(false).Once()
	mockResult.On("Record").Return(mockRecord)
	return mockResult
}

// ResourceWithoutParentGetByID returns a resource fetched when GetById is called
func ResourceWithoutParentGetByID() *Result {
	mockResult := &Result{}
	mockRecord := &Record{}
	mockRecord.On("GetByIndex", 0).Return("test-resource")
	mockRecord.On("GetByIndex", 1).Return("test-source-id")
	mockRecord.On("GetByIndex", 2).Return(nil)
	mockResult.On("Next").Return(true).Once()
	mockResult.On("Next").Return(false).Once()
	mockResult.On("Record").Return(mockRecord)
	return mockResult
}

// ResourceWithParent mocks a resource returned without any parent
func ResourceWithParent() *Result {
	mockResult := &Result{}
	mockRecord := &Record{}
	mockRecord.On("GetByIndex", 0).Return("test-resource")
	mockRecord.On("GetByIndex", 1).Return("test-source-id")
	mockRecord.On("GetByIndex", 2).Return("test-id")
	mockRecord.On("GetByIndex", 3).Return("parent-id")
	mockResult.On("Next").Return(true).Once()
	mockResult.On("Next").Return(false).Once()
	mockResult.On("Record").Return(mockRecord)
	return mockResult
}

// GetEmptyResult represents a case when there are no resources in the database
func GetEmptyResult() *Result {
	mockResult := &Result{}
	mockRecord := &Record{}
	mockResult.On("Next").Return(false).Once()
	mockResult.On("Record").Return(mockRecord)
	return mockResult
}

// CreateResourceSuccessful mimics a condition when a resource was added and uuid is returned
func CreateResourceSuccessful() *Result {
	mockResult := &Result{}
	mockRecord := &Record{}
	mockRecord.On("GetByIndex", 0).Return("test-id")
	mockResult.On("Next").Return(true).Once()
	mockResult.On("Next").Return(false).Once()
	mockResult.On("Record").Return(mockRecord)
	return mockResult
}

// CreateEdgeSuccessful mimics a condition when an edge was created successfully
func CreateEdgeSuccessful() *Result {
	mockResult := &Result{}
	mockRecord := &Record{}
	mockRecord.On("GetByIndex", 0).Return("OWNED_BY")
	mockResult.On("Next").Return(true).Once()
	mockResult.On("Record").Return(mockRecord)
	return mockResult
}

// CreateEdgeFails mimics a condition when an edge was created successfully
func CreateEdgeFails() *Result {
	mockResult := &Result{}
	mockRecord := &Record{}
	mockRecord.On("GetByIndex", 0).Return("")
	mockResult.On("Next").Return(true).Once()
	mockResult.On("Record").Return(mockRecord)
	return mockResult
}

// DeleteEdgeFails mimics a case when the relationships updated is 0
func DeleteEdgeFails() *Result {
	mockCounter := &Counters{}
	mockCounter.On("RelationshipsDeleted").Return(0)
	mockSummary := &ResultSummary{}
	mockSummary.On("Counters").Return(mockCounter)
	mockResult := &Result{}
	mockResult.On("Next").Return(true).Once()
	mockResult.On("Summary").Return(mockSummary, nil)
	return mockResult
}

// SummaryFailure mimics a case when the calling the summary results in failrue
func SummaryFailure() *Result {
	mockResult := &Result{}
	mockResult.On("Next").Return(true).Once()
	mockResult.On("Summary").Return(nil, errors.New("Some summary error"))
	return mockResult
}

// DeleteResourceSuccessful mimics a condition when a resource is deleted succesfully
func DeleteResourceSuccessful() *Result {
	mockCounter := &Counters{}
	mockCounter.On("NodesDeleted").Return(1)
	mockSummary := &ResultSummary{}
	mockSummary.On("Counters").Return(mockCounter)
	mockResult := &Result{}
	mockResult.On("Next").Return(true).Once()
	mockResult.On("Summary").Return(mockSummary, nil)
	return mockResult
}

// DeleteRelationshipSuccessful mimics a condition when a edge is deleted
func DeleteRelationshipSuccessful() *Result {
	mockCounter := &Counters{}
	mockCounter.On("RelationshipsDeleted").Return(1)
	mockSummary := &ResultSummary{}
	mockSummary.On("Counters").Return(mockCounter)
	mockResult := &Result{}
	mockResult.On("Next").Return(true).Once()
	mockResult.On("Summary").Return(mockSummary, nil)
	return mockResult
}

// DeleteResourceNoNodesDeleted mimics a condition when a no nodes were deleted
func DeleteResourceNoNodesDeleted() *Result {
	mockCounter := &Counters{}
	mockCounter.On("NodesDeleted").Return(0)
	mockSummary := &ResultSummary{}
	mockSummary.On("Counters").Return(mockCounter)
	mockResult := &Result{}
	mockResult.On("Next").Return(true).Once()
	mockResult.On("Summary").Return(mockSummary, nil)
	return mockResult
}

// UpdateOwnershipZeroRelationshipsCreated mimics a case when the summary shows zero relationships were created
func UpdateOwnershipZeroRelationshipsCreated() *Result {
	mockCounter := &Counters{}
	mockCounter.On("RelationshipsCreated").Return(0)
	mockSummary := &ResultSummary{}
	mockSummary.On("Counters").Return(mockCounter)
	mockResult := &Result{}
	mockResult.On("Next").Return(true).Once()
	mockResult.On("Summary").Return(mockSummary, nil)
	return mockResult
}

// UpdateOwnershipNoErrors mimics a case when everything went fine - no errors returned
func UpdateOwnershipNoErrors() *Result {
	mockCounter := &Counters{}
	mockCounter.On("RelationshipsCreated").Return(1)
	mockSummary := &ResultSummary{}
	mockSummary.On("Counters").Return(mockCounter)
	mockResult := &Result{}
	mockResult.On("Next").Return(true).Once()
	mockResult.On("Summary").Return(mockSummary, nil)
	return mockResult
}
