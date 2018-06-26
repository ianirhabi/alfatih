// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package docv_test

import (
	"os"
	"testing"
	"time"

	"git.qasico.com/cuxs/docv"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	// run tests
	res := m.Run()

	// cleanup
	os.Exit(res)
}

func TestVersioning(t *testing.T) {
	if er := docv.New(); er != nil {
		panic(er.Error())
	}

	defer docv.Session.Close()

	id := int64(time.Now().Day())
	docType := "document_testing"

	// test cleaning data
	e := docv.Clean(docType, id)
	assert.NoError(t, e)

	data := []struct {
		ID   int64
		Name string
		Doc  *docv.Document
	}{
		{ID: id, Name: "test 1", Doc: &docv.Document{}},
		{ID: id, Name: "test 2", Doc: &docv.Document{}},
		{ID: id, Name: "test 3", Doc: &docv.Document{}},
		{ID: id, Name: "test 4", Doc: &docv.Document{}},
		{ID: id, Name: "test 5", Doc: &docv.Document{}},
		{ID: id, Name: "test 6", Doc: &docv.Document{}},
		{ID: id, Name: "test 7", Doc: &docv.Document{}},
	}

	for i, d := range data {
		// test storing document
		e := docv.Create(docType, id, d, nil)
		assert.NoError(t, e)

		// test get last version document
		cd, e := docv.Show(docType, id)
		assert.NoError(t, e)
		assert.Equal(t, i+1, cd.Version)
	}

	// test get all document version
	vs, e := docv.Get(docType, id)
	assert.NoError(t, e)
	assert.Equal(t, len(data), len(vs))
}
