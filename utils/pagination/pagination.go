package pagination

import "math"

type Data struct {
	currentPage    int64
	maxPages       int64
	resultsCount   int64
	resultsPerPage int64
}

func New(currentPage, resultsCount, resultsPerPage int64) *Data {
	maxPages := int64(1)
	if resultsCount > 0 {
		maxPages = int64(
			math.Ceil(
				float64(resultsCount) / float64(resultsPerPage),
			),
		)
	}

	if currentPage < 1 {
		currentPage = 1
	} else if currentPage > maxPages {
		currentPage = maxPages
	}

	return &Data{
		currentPage:    currentPage,
		maxPages:       maxPages,
		resultsCount:   resultsCount,
		resultsPerPage: resultsPerPage,
	}
}

func (d *Data) CurrentPage() int64 {
	return d.currentPage
}

func (d *Data) MaxPages() int64 {
	return d.maxPages
}

func (d *Data) ResultsCount() int64 {
	return d.resultsCount
}

func (d *Data) ResultsPerPage() int64 {
	return d.resultsPerPage
}

// Returns [LIMIT, OFFSET]
func (d *Data) Limit() (int64, int64) {
	return d.resultsPerPage, d.currentPage*d.resultsPerPage - d.resultsPerPage
}
