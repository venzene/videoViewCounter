package model

type TestCase struct {
	TestName       string
	Vid            string
	ExpectedViews  int
	NParams        int
	TestInput      []VideoInfo
	ExpectedResult []VideoInfo
	ExpectedErr    error
}
