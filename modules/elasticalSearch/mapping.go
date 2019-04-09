package elasticalSearch

//	number_of_shards 数据分片数
//	number_of_replicas 数据备份数,如果只有一台机器，设置为0
func GetMapping() string {
	return `{
	"settings":{
		"number_of_shards":1,
		"number_of_replicas":0,
		"max_result_window":1000
	},
	"doc":{
		"properties":{
			"Id":       	{"type":"integer"},
            "Wid":     		{"type":"integer"},
            "ActivityId":   {"type":"integer"},
            "Alias":        {"type":"text"},
            "ClickKey":     {"type":"text"},
            "Success":      {"type":"keyword","index": "false"},
            "Fail":         {"type":"keyword","index": "false"},
            "NoPrizeReturn":   {"type":"keyword","index": "false"},
            "Extra":        {"type":"keyword","index": "false"},
            "Type":         {"type":"integer"},
            "Disabled":     {"type":"integer"},
            "CreatedAt":    {"type":"date","format":"dateOptionalTime"},
            "UpdatedAt":    {"type":"date","format":"dateOptionalTime"}
		}
	}
}`
}
