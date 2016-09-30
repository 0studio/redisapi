package redisapi

import "testing"
import "fmt"
import "github.com/garyburd/redigo/redis"

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

func TestZ(t *testing.T) {
	client, err := InitDefaultClient(":6379")
	if err != nil {
		t.Errorf("%s\r\n", err.Error())
	}

	err = client.Zadd("key", 100, "mem1")
	if err != nil {
		t.Errorf("%s\r\n", err.Error())
	}
	err = client.Zadd("key", 300, "mem3")
	if err != nil {
		t.Errorf("%s\r\n", err.Error())
	}
	err = client.Zadd("key", 200, "mem2")
	if err != nil {
		t.Errorf("%s\r\n", err.Error())
	}
	cnt, err := client.Zcard("key")
	if err != nil {
		t.Errorf("%s\r\n", err.Error())
	}
	if cnt != 3 {
		t.Errorf("zcard error\r\n")
	}
	rank, err := client.ZRrank("key", "mem3")
	if err != nil {
		t.Errorf("%s\r\n", err.Error())
	}
	if rank != 2 {
		t.Errorf("zrank error\r\n", rank)
	}
	list, err := client.ZRrange("key", 0, -1)
	if len(list) != 3 {
		t.Errorf("zrange error\r\n", rank)
	}

	if list[0].GetMemberAsString() != "mem1" {
		t.Errorf("zrange order error\r\n", rank)
	}
	if list[1].GetMemberAsString() != "mem2" {
		t.Errorf("zrange order error\r\n", rank)
	}
	if list[2].GetMemberAsString() != "mem3" {
		t.Errorf("zrange order error\r\n", rank)
	}

	list, err = client.ZRevRrange("key", 0, -1)
	if list[0].GetMemberAsString() != "mem3" {
		t.Errorf("zrevrange order error\r\n", rank)
	}
	if list[1].GetMemberAsString() != "mem2" {
		t.Errorf("zrevrange order error\r\n", rank)
	}
	if list[2].GetMemberAsString() != "mem1" {
		t.Errorf("zrevrange order error\r\n", rank)
	}

}

func TestZRem(t *testing.T) {
	client, err := InitDefaultClient(":6379")
	if err != nil {
		t.Errorf("%s\r\n", err.Error())
	}

	err = client.Zadd("key", 100, "mem1")
	err = client.Zadd("key", 300, "mem3")
	err = client.Zadd("key", 200, "mem2")
	cnt, err := client.Zcard("key")
	if cnt != 3 {
		t.Errorf("zcard error\r\n")
	}
	err = client.ZRemRangeByRank("key", 0, 1)
	if err != nil {
		t.Errorf("ZRemRangeByRank error\r\n")
	}
	list, err := client.ZRrange("key", 0, -1)
	if len(list) != 1 || list[0].GetMemberAsString() != "mem3" {
		t.Errorf("ZRemRangeByRank error\r\n")
	}

}

func TestZRem2(t *testing.T) {
	client, err := InitDefaultClient(":6379")
	if err != nil {
		t.Errorf("%s\r\n", err.Error())
	}

	err = client.Zadd("key", 100, "mem1")
	err = client.Zadd("key", 300, "mem3")
	err = client.Zadd("key", 200, "mem2")
	cnt, err := client.Zcard("key")
	if cnt != 3 {
		t.Errorf("zcard error\r\n")
	}
	err = client.ZRemRangeByScore("key", "99", "201")
	if err != nil {
		t.Errorf("ZRemRangeByRank error\r\n")
	}
	list, err := client.ZRrange("key", 0, -1)
	if len(list) != 1 || list[0].GetMemberAsString() != "mem3" {
		t.Errorf("ZRemRangeByRank error\r\n")
	}

}

func TestZRem3(t *testing.T) {
	client, err := InitDefaultClient(":6379")
	if err != nil {
		t.Errorf("%s\r\n", err.Error())
	}

	err = client.Zadd("key", 100, "mem1")
	err = client.Zadd("key", 300, "mem3")
	err = client.Zadd("key", 200, "mem2")
	cnt, err := client.Zcard("key")
	if cnt != 3 {
		t.Errorf("zcard error\r\n")
	}
	err = client.ZRemRangeByScore("key", "99", "inf")
	if err != nil {
		t.Errorf("ZRemRangeByRank error\r\n")
	}
	list, err := client.ZRrange("key", 0, -1)
	if len(list) != 0 {
		t.Errorf("ZRemRangeByRank error\r\n", len(list))
	}

}

func TestZRem4(t *testing.T) {
	client, err := InitDefaultClient(":6379")
	if err != nil {
		t.Errorf("%s\r\n", err.Error())
	}

	err = client.Zadd("key", 100, "mem1")
	err = client.Zadd("key", 300, "mem3")
	err = client.Zadd("key", 200, "mem2")
	cnt, err := client.Zcard("key")
	if cnt != 3 {
		t.Errorf("zcard error\r\n")
	}
	err = client.ZRemRangeByScore("key", 99, 201)
	if err != nil {
		t.Errorf("ZRemRangeByRank error\r\n")
	}
	list, err := client.ZRrange("key", 0, -1)
	if len(list) != 1 || list[0].GetMemberAsString() != "mem3" {
		t.Errorf("ZRemRangeByRank error\r\n")
	}

}

func TestZScore(t *testing.T) {
	client, err := InitDefaultClient(":6379")
	if err != nil {
		t.Errorf("%s\r\n", err.Error())
	}
	err = client.Zadd("key", 100, "mem1")
	err = client.Zadd("key", 10000000000000000, "mem3")
	err = client.Zadd("key", 200, "mem2")
	score, err := client.ZScore("key", "mem1")
	if err != nil {
		t.Error("zscore error", err)
	}
	if score == 0 {
		t.Error("zscore error", err)
	}

	score, err = client.ZScore("key", "mem3")
	if err != nil {
		t.Error("zscore error", err)
	}
	if score == 0 {
		t.Error("zscore error", err)
	}

	score, err = client.ZScore("key", "not_exists")
	if err != redis.ErrNil {
		t.Error("zscore should be nil")
	}

	fmt.Println(score, err)

}

func TestZScoreFloat64(t *testing.T) {
	client, err := InitDefaultClient(":6379")
	if err != nil {
		t.Errorf("%s\r\n", err.Error())
	}
	err = client.Zadd("key", 100, "mem1")
	err = client.Zadd("key", 10000000000000000, "mem3")
	err = client.Zadd("key", 1E+100, "mem2")
	score, err := client.ZScoreAsFloat64("key", "mem1")
	if err != nil {
		t.Error("zscore error", err)
	}
	if score == 0 {
		t.Error("zscore error", err)
	}

	score, err = client.ZScoreAsFloat64("key", "mem3")
	if err != nil {
		t.Error("zscore error", err)
	}
	if score == 0 {
		t.Error("zscore error", err)
	}
	fmt.Println(score)

	score, err = client.ZScoreAsFloat64("key", "mem2")
	if err != nil {
		t.Error("zscore error", err)
	}
	if score == 0 {
		t.Error("zscore error", err)
	}
	fmt.Println(score)

}

func TestSadd(t *testing.T) {
	client, err := InitDefaultClient(":6379")
	if err != nil {
		t.Errorf("%s\r\n", err.Error())
	}
	err = client.Sadd("setkey", "mem1")
	if err != nil {
		t.Errorf("%s\r\n", err.Error())
	}
	err = client.Sadd("setkey", "mem2", "mem3")
	n, err := client.Scard("setkey")
	if n != 3 {
		t.Errorf("scard \r\n", err)

	}
	exits := client.SisMember("setkey", "mem1")
	if !exits {
		t.Errorf("SisMembererror \r\n", err)

	}
	members, err := client.Smembers("setkey")
	fmt.Println(members)
	if len(members) != 3 {
		t.Errorf("Smembers \r\n", err)

	}
	members2, err := client.SmembersAsString("setkey")
	fmt.Println("members", members2)
	if len(members2) != 3 {
		t.Errorf("Smembers \r\n", err)
	}

	err = client.Srem("setkey", "mem1")
	if err != nil {
		t.Errorf("Srem \r\n", err)
	}
	m, err := client.SpopAsString("setkey")
	fmt.Println(m)
	if err != nil {
		t.Errorf("SpopAsString \r\n", err)
	}
	m3, err := client.SrandMemberAsString("setkey")
	fmt.Println(m3)
	if err != nil {
		t.Errorf("SrandMember \r\n", err)
	}

}


func TestGeo(t *testing.T){
	client, err := InitDefaultClient(":6379")
	if err != nil {
		t.Errorf("%s\r\n", err.Error())
	}
	err = client.GeoAdd("test",Coordinate{Latitude:22.8777898943,Longitude:114.6834120888},"huizhou")
	if err != nil{
		t.Errorf("%s\r\n", err.Error())
	}

	err = client.GeoAdd("test",Coordinate{Latitude:22.4421753709,Longitude:114.1678601427},"hongkong")
	if err != nil{
		t.Errorf("%s\r\n", err.Error())
	}
	BaoanCoordinate := Coordinate{Latitude:22.6679807509,Longitude:113.8106668351}
	err = client.GeoAdd("test",BaoanCoordinate,"baoan")
	if err != nil{
		t.Errorf("%s\r\n", err.Error())
	}

	v,err := client.GeoPos("test","baoan")
	if err != nil{
		t.Errorf("%s\r\n", err.Error())
	}
	fmt.Print(v)
	dis, err := client.GeoDist("test","baoan","hongkong",DistanceUnitKM)
	if err != nil{
		t.Errorf("%s\r\n", err.Error())
	}
	fmt.Println(dis)


	c,err := client.GeoPos("test","baoan")
	if err != nil{
		t.Errorf("%s\r\n", err.Error())
	}
	fmt.Println(c)


	coor,err := client.GeoRadius("test",BaoanCoordinate,2000,SetDistancUnit(DistanceUnitKM),SetSort(ASC),SetWith(WITHBOth))
	if err != nil{
		t.Errorf("%s\r\n", err.Error())
	}
	fmt.Println(coor)


	coor,err = client.GeoRadiusByMember("test","hongkong",200,SetDistancUnit(DistanceUnitKM),SetSort(ASC),SetWith(WITHBOth))
	if err != nil{
		t.Errorf("%s\r\n", err.Error())
	}
	fmt.Printf("%s\n",(*coor)[2].value)

}