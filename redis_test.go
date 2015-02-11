package redisapi

import "testing"

func TestExists(t *testing.T) {
	client, err := InitDefaultClient(":6379")
	if err != nil {
		t.Errorf("%s\r\n", err.Error())
	}

	client.Set("aaa", []byte("bbb"))
	result := client.Exists("aaa")
	if result != true {
		t.Errorf("result should be true")
	}

	result = client.Exists("bbb")
	if result != false {
		t.Errorf("result should be false")
	}

	client.Lpush("aaaa", "bbbb")
	v, err := client.Brpop("aaaa", -1)
	value := v.(string)
	if value != "bbbb" {
		t.Errorf("result should be bbbb")
	}

	client.ClearAll()
	client.Set("aaa", []byte("bbb"))
	v, err = client.Get("aaa")
	if string(v.([]byte)) != "bbb" {
		t.Errorf("result should be bbb")
	}

	client.ClearAll()
	client.Incr("aaa", 1)
	v, err = client.Get("aaa")
	if string(v.([]byte)) != "1" {
		t.Errorf("result should be 1")
	}
	client.ClearAll()

	kvMap := make(map[string][]byte)
	kvMap["1"] = []byte("123")
	kvMap["2"] = []byte("456")
	client.MultiSet(kvMap)

	v, err = client.Get("1")
	if string(v.([]byte)) != "123" {
		t.Errorf("result should be 123")
	}

	v, err = client.MultiGet([]interface{}{"1", "2"})
	if string(v.([]interface{})[0].([]byte)) != "123" {
		t.Errorf("result should be 123")
	}
	if string(v.([]interface{})[1].([]byte)) != "456" {
		t.Errorf("result should be 456")
	}

	client.Set("aaa", []byte("bbb"))
	client.Delete("aaa")
	v, err = client.Get("aaa")
	if len(v.([]byte)) != 0 {
		t.Errorf("result len should be 0")
	}
	client.ClearAll()

	for i := 0; i < 100; i++ {
		client.Rpush("aaa", i)
	}
	reli, err := client.Lrange("aaa", 0, 19)

	if len(reli) != 20 {
		t.Errorf("result len should be 20")
	}
	client.ClearAll()

}

func TestPubSub(t *testing.T) {
	client, err := InitDefaultClient(":6379")
	if err != nil {
		t.Errorf("%s\r\n", err.Error())
	}

	err = client.Pub("first", "123")
	if err != nil {
		t.Errorf("%s\r\n", err.Error())
	}

	v, err := client.Sub("first", "123")
	if err != nil {
		t.Errorf("%s\r\n", err.Error())
	}

	if v[0] != "first" {
		t.Errorf("v[0] should be first , but is %s", v[0])
	}
	if v[1] != "123" {
		t.Errorf("v[1] should be 123 , but is %s", v[1])
	}

}

func TestHash(t *testing.T) {
	client, err := InitDefaultClient(":6379")
	if err != nil {
		t.Errorf("%s\r\n", err.Error())
	}

	table := "hash_test"
	var scoreList []ScoreStruct
	scoreList = append(scoreList, ScoreStruct{Member: "111111", Score: 111111})
	scoreList = append(scoreList, ScoreStruct{Member: "222222", Score: 222222})
	err = client.HMset(table, scoreList)
	if err != nil {
		t.Errorf("%s\r\n", err.Error())
	}

	scoreList, err = client.HMget(table, "111111", "222222")
	if err != nil {
		t.Errorf("%s\r\n", err.Error())
	}
	t.Log(scoreList)
}
