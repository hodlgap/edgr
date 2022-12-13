# edgr | Makes SEC filings not terrible

## `github.com/edgr/core` module

### Installation
Supports Go v1.19+ and modules. Run the usual
```sh
go get -u github.com/hodlgap/edgr
```

### Usage
There are 3 easy functions
```go
// GetPublicCompanies returns a list of public companies.
func GetPublicCompanies() ([]Company, error) {...}

// GetFiler gets a single filer from the SEC website based on symbol.
func GetFiler(symbol string) (filer *model.Filer, err error) {...}

// GetFilings gets a list of filings for a single CIK.
func GetFilings(cik, formtype, stoptime string) (filings []SECFiling, err error) {...}
```
