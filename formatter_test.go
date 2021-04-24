package logrus

import (
	//"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrefixFieldClashesFuncReportCaller(t *testing.T) {
	data := Fields{FieldKeyFunc: "CustomFunc", FieldKeyFile: "custom.go"}
	var fieldMap FieldMap
	var value interface{}
	var hasKey bool

	prefixFieldClashes(data, fieldMap, true)

	_, hasKey = data[FieldKeyFunc]
	assert.False(t, hasKey, "func not deleted when ReportCaller=true")

	_, hasKey = data[FieldKeyFile]
	assert.False(t, hasKey, "file not deleted when ReportCaller=true")

	value, hasKey = data["fields."+FieldKeyFunc]
	assert.True(t, hasKey, "fields.func not set when ReportCaller=true")
	assert.Equal(t, "CustomFunc", value, "fields.func not set as expected when ReportCaller=true")

	value, hasKey = data["fields."+FieldKeyFile]
	assert.True(t, hasKey, "fields.file not set when ReportCaller=true")
	assert.Equal(t, "custom.go", value, "fields.file not set as expected when ReportCaller=true")
}

func TestPrefixFieldClashesFuncNoReportCaller(t *testing.T) {
	data := Fields{FieldKeyFunc: "CustomFunc", FieldKeyFile: "custom.go"}
	var fieldMap FieldMap
	var value interface{}
	var hasKey bool

	prefixFieldClashes(data, fieldMap, false)

	value, hasKey = data[FieldKeyFunc]
	assert.True(t, hasKey, "func deleted when ReportCaller=false")
	assert.Equal(t, "CustomFunc", value, "func set when ReportCaller=false")

	value, hasKey = data[FieldKeyFile]
	assert.True(t, hasKey, "file deleted when ReportCaller=false")
	assert.Equal(t, "custom.go", value, "file set when ReportCaller=false")

	value, hasKey = data["fields."+FieldKeyFunc]
	assert.False(t, hasKey, "fields.func set when ReportCaller=false")

	value, hasKey = data["fields."+FieldKeyFile]
	assert.False(t, hasKey, "fields.file set when ReportCaller=false")
}
