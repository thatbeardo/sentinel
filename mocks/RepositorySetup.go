package mocks

// ReturnResultWithData mimics a case when everything went fine - no errors returned
func ReturnResultWithData(data interface{}) *Result {
	mockResult := &Result{}
	mockResult.On("Next").Return(true).Once()
	mockResult.On("Next").Return(false).Once()
	mockRecord := &Record{}
	mockRecord.On("GetByIndex", 0).Return(data)
	mockResult.On("Record").Return(mockRecord)
	return mockResult
}
