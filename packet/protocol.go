package packet

type Body map[string]interface{}

// implement error interface
func (b Body) Error() string {
	return b["msg"].(string)
}

// func Decode(e interface{}, body []byte) {
// 	net := bytes.NewBuffer(body)
//
// 	decoder := gob.NewDecoder(net)
//
// 	err := decoder.Decode(e)
// 	if err != nil {
// 		log.Println(err.Error())
// 		return
// 	}
// 	log.Println("%+v", e)
// }
