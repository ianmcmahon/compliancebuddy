package faaservices

const (
	FAARegistryDBDownloadURL string = "http://registry.faa.gov/database/ReleasableAircraft.zip"
)

type Client struct {
	RegistryFetcher RegistryFetcher
}

func NewClient() *Client {
	return &Client{
		RegistryFetcher: LiveRegistryFetcher{FAARegistryDBDownloadURL},
	}
}

//////

// type AD struct {
// 	DocumentNumber string
// 	Url            url.URL
// 	Subject        string
// 	Title          string
// 	Type           string
// 	RegulatoryText string
// }

// type count int

// func adCount(product model.Product) (int, error) {
// 	client := &http.Client{}

// 	model := product.GetModel()

// 	req, err := http.NewRequest(http.MethodGet, "http://services.faa.gov/document/ad/count", nil)
// 	if err != nil {
// 		return 0, err
// 	}

// 	q := req.URL.Query()
// 	q.Add("format", "application/xml")
// 	q.Add("make", model.Make)
// 	q.Add("model", model.Model)
// 	req.URL.RawQuery = q.Encode()

// 	resp, err := client.Do(req)
// 	defer resp.Body.Close()
// 	if err != nil {
// 		return 0, err
// 	}

// 	var countResp count = -1

// 	dec := xml.NewDecoder(resp.Body)
// 	if err := dec.Decode(&countResp); err != nil {
// 		return 0, err
// 	}

// 	return int(countResp), nil
// }

// type document struct {
// 	DocumentNumber string `xml:"Documentnumber"`
// 	Title          string `xml:"Title"`
// 	Type           string `xml:"Type"`
// 	Uri            string `xml:"Uri"`
// 	Subject        string `xml:"Subject"`
// }

// type documents struct {
// 	List []document `xml:"Document"`
// }

// func adList(product model.Product, offset int) ([]document, error) {
// 	client := &http.Client{}

// 	model := product.GetModel()

// 	req, err := http.NewRequest(http.MethodGet, "http://services.faa.gov/document/ad/list", nil)
// 	if err != nil {
// 		return nil, err
// 	}

// 	q := req.URL.Query()
// 	q.Add("format", "application/xml")

// 	q.Add("fields", "Documentnumber")
// 	q.Add("fields", "Title")
// 	q.Add("fields", "Type")
// 	q.Add("fields", "Uri")
// 	q.Add("fields", "Subject")

// 	q.Add("sort", "effectivedate")

// 	q.Add("limit", "100")
// 	q.Add("offset", fmt.Sprintf("%d", offset))

// 	q.Add("make", model.Make)
// 	q.Add("model", model.Model)

// 	req.URL.RawQuery = q.Encode()

// 	resp, err := client.Do(req)
// 	defer resp.Body.Close()
// 	if err != nil {
// 		return nil, err
// 	}

// 	docs := documents{}

// 	dec := xml.NewDecoder(resp.Body)
// 	if err := dec.Decode(&docs); err != nil {
// 		return nil, err
// 	}

// 	return docs.List, nil
// }

// func ADSearch(product model.Product) ([]AD, error) {
// 	// have to call count service first to get total count for the query, then paginate based on that number
// 	// FAA services are very badly written.
// 	// they support json, but the json they return is invalid

// 	count, err := adCount(product)
// 	if err != nil {
// 		return nil, err
// 	}

// 	fmt.Printf("count: %d\n", count)

// 	seen := 0

// 	docs := []document{}

// 	for seen < count {
// 		page, err := adList(product, seen+1)
// 		if err != nil {
// 			return nil, err
// 		}
// 		docs = append(docs, page...)

// 		seen += len(page)
// 	}

// 	fmt.Printf("%d docs retrieved\n", len(docs))

// 	for _, doc := range docs {
// 		fmt.Printf("\t%s: %s\n", doc.DocumentNumber, doc.Subject)
// 	}

// 	return []AD{}, nil
// }
