package models

type Goods struct {
	Id            int     `json:"id"`
	Title         string  `json:"title"`
	SubTitle      string  `json:"sub_title"`
	GoodsSn       string  `json:"goods_sn"`
	CateId        int     `json:"cate_id"`
	ClickCount    int     `json:"click_count"`
	GoodsNumber   int     `json:"goods_number"`
	Price         float64 `json:"price"`
	MarketPrice   float64 `json:"market_price"`
	RelationGoods string  `json:"relation_goods"`
	GoodsAttr     string  `json:"goods_attr"`
	GoodsVersion  string  `json:"goods_version"`
	GoodsImg      string  `json:"goods_img"`
	GoodsGift     string  `json:"goods_gift"`
	GoodsFitting  string  `json:"goods_fitting"`
	GoodsColor    string  `json:"goods_color"`
	GoodsKeywords string  `json:"goods_keywords"`
	GoodsDesc     string  `json:"goods_desc"`
	GoodsContent  string  `json:"goods_content"`
	IsDelete      int     `json:"is_delete"`
	IsHot         int     `json:"is_hot"`
	IsBest        int     `json:"is_best"`
	IsNew         int     `json:"is_new"`
	GoodsTypeId   int     `json:"goods_type_id"`
	Sort          int     `json:"sort"`
	Status        int     `json:"status"`
	AddTime       int     `json:"add_time"`
	MerchantId    int     `json:"merchant_id"` // 商户ID，0表示平台商品
}

func (Goods) TableName() string {
	return "goods"
}

/*
根据商品分类获取推荐商品
	@param {Number} cateId - 分类id
	@param {String} goodsType -  hot  best  new all
	@param {Number} limitNum -  数量

	1  表示顶级分类
		21
		23
		24


*/

func GetGoodsByCategory(cateId int, goodsType string, limitNum int) []Goods {

	//判断cateId 是否是顶级分类
	goodsCate := GoodsCate{Id: cateId}
	DB.Find(&goodsCate)
	var tempSlice []int
	if goodsCate.Pid == 0 { //顶级分类
		//获取顶级分类下面的二级分类
		goodsCateList := []GoodsCate{}
		DB.Where("pid = ?", goodsCate.Id).Find(&goodsCateList)

		for i := 0; i < len(goodsCateList); i++ {
			tempSlice = append(tempSlice, goodsCateList[i].Id)
		}

	}
	tempSlice = append(tempSlice, cateId)

	goodsList := []Goods{}
	where := "cate_id in ?"
	switch goodsType {
	case "hot":
		where += " AND is_hot=1"
	case "best":
		where += " AND is_best=1"
	case "new":
		where += " AND is_new=1"
	default:
		break
	}

	DB.Where(where, tempSlice).Select("id,title,price,goods_img,sub_title").Limit(limitNum).Order("sort desc").Find(&goodsList)
	return goodsList
}
