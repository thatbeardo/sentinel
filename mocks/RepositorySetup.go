package mocks

// GetMockResult returns
func GetMockResult() *Result {
	mockResult := &Result{}
	mockRecord := &Record{}
	mockRecord.On("GetByIndex", 0).Return("test-resource")
	mockRecord.On("GetByIndex", 1).Return("test-source-id")
	mockRecord.On("GetByIndex", 2).Return("test-id")
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
