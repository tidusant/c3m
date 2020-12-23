package main

//repush
import (
	"context"
	"fmt"
	"github.com/tidusant/c3m/common/c3mcommon"
	"github.com/tidusant/c3m/common/mycrypto"
	"github.com/tidusant/c3m/common/mystring"
	"google.golang.org/grpc"
	"os"
	"time"

	"github.com/tidusant/c3m/common/log"
	"go.mongodb.org/mongo-driver/bson/primitive"

	pb "github.com/tidusant/c3m/grpc/protoc"
	rpch "github.com/tidusant/c3m/repo/cuahang"
	"github.com/tidusant/c3m/repo/models"

	"encoding/base64"
	"encoding/json"
	"math"

	"net"

	"strconv"
	"strings"
)

const (
	name string = "ord"
	ver  string = "1"
)

type service struct {
	pb.UnimplementedGRPCServicesServer
}

func (s *service) Call(ctx context.Context, in *pb.RPCRequest) (*pb.RPCResponse, error) {
	resp := &pb.RPCResponse{Data: "Hello " + in.GetAppName(), RPCName: name, Version: ver}
	rs := models.RequestResult{Error: ""}
	//get input data into user session
	var usex models.UserSession
	usex.Session = in.Session
	usex.Action = in.Action

	usex.UserIP = in.UserIP
	usex.Params = in.Params
	usex.UserID, _ = primitive.ObjectIDFromHex(in.UserID)

	log.Debugf("shopid:%v", in.ShopID)
	//check shop permission
	if in.ShopID != "" {
		shopidObj, _ := primitive.ObjectIDFromHex(in.ShopID)
		shop := rpch.GetShopById(usex.UserID, shopidObj)
		if shop.Status == 0 {
			rs.Error = "Site is disable"

		}
		usex.Shop = shop
	}
	if rs.Error == "" {
		if usex.Action == "statusc" {
			rs = LoadAllStatusCount(usex)
		} else if usex.Action == "status" {
			rs = LoadAllStatus(usex)
		} else if usex.Action == "lao" {
			rs = LoadAllOrderByStatus(usex)
		} else if usex.Action == "lg" {
			rs = LoadCities(usex)
		} else if usex.Action == "us" {
			rs = UpdateOrderStatus(usex)
		} else if usex.Action == "ds" {
			rs = DeleteOrderStatus(usex)
		} else if usex.Action == "ss" {
			rs = SaveStatus(usex)
		} else if usex.Action == "uo" {
			rs = UpdateOrder(usex)
		}
	}
	//convert RequestResult into json
	log.Debugf("data return:%+v", rs)
	b, _ := json.Marshal(rs)
	resp.Data = string(b)
	return resp, nil
}

//parse order from web (status="") and update
func parseOrder(order models.Order, usex models.UserSession, defaultstatus models.OrderStatus) {
	if order.C == "" {
		return
	}
	//loop to get item
	orderc := order.C

	for {
		if len(orderc) < 3 {
			break
		}

		code := orderc[:3]
		orderc = orderc[3:]
		//get prod
		prod := rpch.GetProdByCode(usex.Shop.ID, code)
		if prod.Code == "" {
			//prod not found
			break
		}

		//get num
		num := 1
		//loop to get num
		numstr := ""
		for {
			if len(orderc) <= 0 {
				break
			}
			str := orderc[:1]
			if !mystring.IsInt(str) {
				break
			}
			numstr = numstr + str
			orderc = orderc[1:]
		}

		//check num
		if numstr != "" {
			num, _ = strconv.Atoi(numstr)
		}

		//create order item repush
		var item models.OrderItem
		item.Code = prod.Code
		item.BasePrice = prod.Langs[order.L].BasePrice
		item.Price = prod.Langs[order.L].Price
		item.Title = prod.Langs[order.L].Name
		item.Avatar = prod.Langs[order.L].Avatar
		item.Num = num
		order.Items = append(order.Items, item)
		order.Total += item.Price
		order.BaseTotal += item.BasePrice

		order.PartnerShipFee = order.ShipFee

	}
	order.Status = defaultstatus.ID.Hex()
	rpch.SaveOrder(order)
}
func LoadAllOrderByStatus(usex models.UserSession) models.RequestResult {

	args := strings.Split(usex.Params, ",")
	status := args[0]
	page := 1
	if len(args) > 1 {
		page, _ = strconv.Atoi(args[1])
	}
	pagesize := int64(100)
	count := 1
	searchterm := ""
	if len(args) > 2 {
		searchterm = args[2]
		if searchterm != "" {
			byteDecode, _ := base64.StdEncoding.DecodeString(mycrypto.Base64fix(searchterm))
			searchterm = string(byteDecode)
			log.Debugf("searchterm: %s", searchterm)
		}
	}

	count = int(rpch.CountOrdersByStatus(usex.Shop.ID, status, searchterm))
	totalPage := (int)(math.Ceil(float64(count) / float64(pagesize)))
	if page > totalPage {
		page = totalPage
	}

	//update order from web

	//orders := rpch.GetOrdersByStatus(usex.Shop.ID, "all", 0, pagesize, "")
	//default status
	defaultstatus := rpch.GetDefaultOrderStatus(usex.Shop.ID)
	if defaultstatus.ID == primitive.NilObjectID {
		return models.RequestResult{Status: 0, Error: "No Default Status for this shop"}
	}
	//default shipper
	//defaultshipper := rpch.GetDefaultShipper(usex.Shop.ID)
	//all campaign
	camps := rpch.GetAllCampaigns(usex.Shop.ID)
	mapcamp := make(map[string]string)
	for _, v := range camps {
		mapcamp[v.ID.Hex()] = v.Name
	}
	//for _, order := range orders {
	//	parseOrder(order, usex, defaultstatus)
	//}

	orders := rpch.GetOrdersByStatus(usex.Shop.ID, status, page, pagesize, searchterm)
	cuss := make(map[string]models.Customer)
	for k, v := range orders {
		//get cus
		if _, ok := cuss[v.Phone]; !ok {
			cuss[v.Phone] = rpch.GetCusByPhone(v.Phone, usex.Shop.ID.Hex())
		}
		orders[k].Name = cuss[v.Phone].Name
		if campname, ok := mapcamp[orders[k].CampaignId]; ok {
			orders[k].CampaignName = campname
		}
		//log.Debugf("after GetCusByPhone %+v",v)
		orders[k].Email = cuss[v.Phone].Email
		orders[k].City = cuss[v.Phone].City
		orders[k].District = cuss[v.Phone].District
		orders[k].Ward = cuss[v.Phone].Ward
		orders[k].Address = cuss[v.Phone].Address
		orders[k].CusNote = cuss[v.Phone].Note
		orders[k].OrderCount = rpch.CountOrderByCus(v.Phone, usex.Shop.ID.Hex())
		orders[k].SearchIndex = ""

	}
	info, _ := json.Marshal(orders)
	strrt := `{"rs":` + string(info) + `,"pagecount":` + strconv.Itoa(totalPage) + `}`
	//strrt = string(info)
	return models.RequestResult{Status: 1, Error: "", Data: strrt}
}
func LoadAllStatus(usex models.UserSession) models.RequestResult {

	//default status
	status := rpch.GetAllOrderStatus(usex.Shop.ID)

	info, _ := json.Marshal(status)

	strrt := string(info)
	return models.RequestResult{Status: 1, Error: "", Data: strrt}
}
func LoadAllStatusCount(usex models.UserSession) models.RequestResult {

	//default status
	status := rpch.GetAllOrderStatus(usex.Shop.ID)
	type orderStat struct {
		Id      string `bson:"id"`
		Name    string `bson:"name"`
		Count   int    `bson:"count"`
		Color   string
		Default bool
		Finish  bool
	}
	var lstOrderStat []orderStat
	defaultstat := ""
	for k, v := range status {
		status[k].OrderCount = int(rpch.CountOrdersByStatus(usex.Shop.ID, v.ID.Hex(), ""))
		if status[k].Default {
			defaultstat = status[k].ID.Hex()
		}
		lstOrderStat = append(lstOrderStat, orderStat{
			Id:      status[k].ID.Hex(),
			Name:    status[k].Title,
			Count:   status[k].OrderCount,
			Color:   status[k].Color,
			Default: status[k].Default,
			Finish:  status[k].Finish,
		})
	}

	info, err := json.Marshal(lstOrderStat)
	if err != nil {
		fmt.Println(err.Error())
	}
	strrt := string(info)
	return models.RequestResult{Status: 1, Error: "", Data: `{"default":"` + defaultstat + `","status":` + strrt + `}`}

}
func LoadCities(usex models.UserSession) models.RequestResult {
	//
	////default status
	//source:=usex.Params
	//cities := rpch.GetCities(source)
	////convert data to have field code for select box
	//strrt:=""
	//if source=="ghtk"{
	//	data:=make(map[string]models.City)
	//	for _,city :=range cities{
	//		var datacity models.City
	//		datacity=city
	//		datacity.Code=datacity.GhtkCode
	//		datacity.Districts=make(map[string]models.District)
	//		for _,district :=range city.Districts{
	//			var datadistrict models.District
	//			datadistrict=district
	//			datadistrict.Code=datadistrict.GhtkCode
	//			datadistrict.Wards=make(map[string]models.Ward)
	//			for _,ward:=range district.Wards{
	//				var dataward models.Ward
	//				dataward=ward
	//				dataward.Code=dataward.GhtkCode
	//				datadistrict.Wards[ward.GhtkCode]=dataward
	//			}
	//			datacity.Districts[district.GhtkCode]=datadistrict
	//		}
	//
	//		data[datacity.GhtkCode]=datacity
	//	}
	//	citiesb, _ := json.Marshal(data)
	//	strrt = string(citiesb)
	//}
	//
	return models.RequestResult{Status: 1, Error: "", Data: ""}
}

func UpdateOrderStatus(usex models.UserSession) models.RequestResult {

	info := strings.Split(usex.Params, ",")
	cancelPartner := "0"
	if len(info) > 1 {
		changestatusid, _ := primitive.ObjectIDFromHex(info[1])
		orderid, _ := primitive.ObjectIDFromHex(info[0])

		//check cancel ghtk status:
		status := rpch.GetStatusByID(changestatusid, usex.Shop.ID)
		ghtkstatussync := status.PartnerStatus["ghtk"]
		if ghtkstatussync != nil {
			for _, statcode := range ghtkstatussync {
				if statcode == "-1" {
					cancelPartner = "1"
				}
			}
		}
		//check stock
		order := rpch.GetOrderByID(orderid, usex.Shop.ID)
		if order.Status == "" {
			return models.RequestResult{Status: 0, Error: "order not found", Data: ""}
		}
		statusid, _ := primitive.ObjectIDFromHex(order.Status)
		oldstat := rpch.GetStatusByID(statusid, usex.Shop.ID)
		if oldstat.Export != status.Export {
			//update stock
			sign := 1
			if oldstat.Export {
				sign = -1 //return to stock
			}
			var exportitems []models.ExportItem
			for _, v := range order.Items {
				prod := rpch.GetProdByCode(usex.Shop.ID, v.ProdCode)
				var newexportitem models.ExportItem
				for _, prop := range prod.Properties {
					if prop.Code == v.Code {
						prop.Stock -= v.Num * sign
						if prop.Stock < 0 {
							titleb, _ := mycrypto.DecompressFromBase64(v.Title)
							return models.RequestResult{Status: 0, Error: "Out of stock: " + string(titleb), Data: ""}

						}
						newexportitem.Code = v.ProdCode
						newexportitem.ItemCode = v.Code
						newexportitem.Num = v.Num * sign
						newexportitem.ShopId = usex.Shop.ID.Hex()
						break
					}
				}
				exportitems = append(exportitems, newexportitem)

			}
			if !rpch.ExportItem(exportitems) {
				return models.RequestResult{Status: 0, Error: "Error update stock", Data: ""}

			}
		}
		rpch.UpdateOrderStatus(usex.Shop.ID, changestatusid, orderid.Hex())

	}
	return models.RequestResult{Status: 1, Error: "", Message: cancelPartner, Data: ""}

}

func SaveStatus(usex models.UserSession) models.RequestResult {

	var status models.OrderStatus
	err := json.Unmarshal([]byte(usex.Params), &status)
	if !c3mcommon.CheckError("update status parse json", err) {
		return models.RequestResult{Status: 0, Error: "update status fail", Data: ""}

	}
	//check old status
	oldstat := status
	if oldstat.ID != primitive.NilObjectID {
		oldstat = rpch.GetStatusByID(status.ID, usex.Shop.ID)
		oldstat.Title = status.Title
		oldstat.Color = status.Color
		oldstat.Default = status.Default
		oldstat.Finish = status.Finish
		oldstat.Export = status.Export
		oldstat.PartnerStatus = status.PartnerStatus
	} else {
		oldstat.UserId = usex.UserID.Hex()
		oldstat.ShopId = usex.Shop.ID.Hex()
	}

	//check default
	if oldstat.Default == true {
		rpch.UnSetStatusDefault(usex.Shop.ID)
	}
	if oldstat.Color == "" {
		oldstat.Color = "ffffff"
	}

	oldstat = rpch.SaveOrderStatus(oldstat)
	b, _ := json.Marshal(oldstat)
	return models.RequestResult{Status: 1, Error: "", Data: string(b)}

}

func DeleteOrderStatus(usex models.UserSession) models.RequestResult {
	//get stat
	statusid, _ := primitive.ObjectIDFromHex(usex.Params)
	stat := rpch.GetStatusByID(statusid, usex.Shop.ID)
	if stat.ID == primitive.NilObjectID {
		return models.RequestResult{Status: 0, Error: "Status not found"}

	}
	if stat.Default {
		return models.RequestResult{Status: 0, Error: "Status is default. Please select another status default."}

	}
	//check status empty
	count := rpch.GetCountOrderByStatus(stat)
	//check old status
	if count > 0 {
		return models.RequestResult{Status: 0, Error: "Status not empty. " + strconv.Itoa(int(count)) + " orders use this status"}

	}

	rpch.DeleteOrderStatus(stat)
	return models.RequestResult{Status: 1, Error: ""}
}

func UpdateOrder(usex models.UserSession) models.RequestResult {
	shop := rpch.GetShopById(usex.UserID, usex.Shop.ID)
	if shop.Status == 0 {
		return models.RequestResult{Status: 0, Error: "Shop is disabled"}

	}
	var order models.Order
	err := json.Unmarshal([]byte(usex.Params), &order)
	if !c3mcommon.CheckError("update order parse json", err) {
		return models.RequestResult{Status: 0, Error: "update order fail"}

	}
	oldorder := order
	mapolditems := make(map[string]models.OrderItem)
	if order.ID != primitive.NilObjectID {
		oldorder = rpch.GetOrderByID(order.ID, shop.ID)
		for _, v := range oldorder.Items {
			mapolditems[v.Code] = v
		}
	} else {
		oldorder.ShopId = usex.Shop.ID.Hex()
		oldorder.ID = primitive.NewObjectID()
		oldorder.Created = time.Now().Unix()
		oldorder.Modified = oldorder.Created
	}

	//all campaign
	camps := rpch.GetAllCampaigns(usex.Shop.ID)
	mapcamp := make(map[string]string)
	for _, v := range camps {
		mapcamp[v.ID.Hex()] = v.Name
	}
	//all shipper
	//shippers := rpch.GetAllShipper(usex.Shop.ID)
	//mapshipper := make(map[string]string)
	//for _, v := range shippers {
	//	mapshipper[v.ID] = v.Name
	//}
	stats := rpch.GetAllOrderStatus(usex.Shop.ID)
	mapstat := make(map[string]models.OrderStatus)
	for _, v := range stats {
		mapstat[v.ID.Hex()] = v
	}

	//check export status
	if mapstat[order.Status].Export {
		//check stock

		var exportitems []models.ExportItem
		for _, v := range order.Items {
			olditemcount := 0
			if _, ok := mapolditems[v.Code]; ok {
				olditemcount = mapolditems[v.Code].Num
			}
			var newexportitem models.ExportItem
			prod := rpch.GetProdByCode(usex.Shop.ID, v.ProdCode)
			for _, prop := range prod.Properties {
				if prop.Code == v.Code {
					prop.Stock -= v.Num - olditemcount
					if prop.Stock < 0 {
						titleb, _ := mycrypto.DecompressFromBase64(v.Title)
						return models.RequestResult{Status: 0, Error: "Out of stock: " + string(titleb)}

					}
					newexportitem.Code = v.ProdCode
					newexportitem.ItemCode = v.Code
					newexportitem.Num = v.Num

					newexportitem.ShopId = usex.Shop.ID.Hex()
					break
				}
			}
			exportitems = append(exportitems, newexportitem)

		}

		//return item stock
		var oldexportitems []models.ExportItem
		for _, v := range mapolditems {
			var old models.ExportItem
			old.ShopId = usex.Shop.ID.Hex()
			old.Code = v.ProdCode
			old.ItemCode = v.Code
			old.Num = -v.Num
			oldexportitems = append(oldexportitems, old)
		}
		rpch.ExportItem(oldexportitems)
		if !rpch.ExportItem(exportitems) {
			return models.RequestResult{Status: 0, Error: "Error update stock"}

		}

	}
	var cus models.Customer
	if oldorder.Phone == order.Phone {
		cus = rpch.GetCusByPhone(order.Phone, shop.ID.Hex())
	} else if oldorder.Phone == "" && oldorder.Email == order.Email {
		cus = rpch.GetCusByEmail(order.Email, shop.ID.Hex())
	}
	cus.Phone = order.Phone
	cus.Name = order.Name
	cus.City = order.City
	cus.District = order.District
	cus.Ward = order.Ward
	cus.Address = order.Address
	cus.Email = order.Email
	cus.Note = order.CusNote
	cus.ShopId = shop.ID.Hex()
	if rpch.SaveCus(cus) {
		//save order
		oldorder.City = order.City
		oldorder.District = order.District
		oldorder.Ward = order.Ward
		oldorder.Address = order.Address
		oldorder.Name = order.Name
		oldorder.Email = order.Email
		oldorder.Phone = order.Phone
		oldorder.CusNote = order.CusNote

		oldorder.Items = order.Items
		oldorder.BaseTotal = order.BaseTotal
		if order.CampaignId != "" {
			oldorder.CampaignId = order.CampaignId
			oldorder.CampaignName = mapcamp[order.CampaignId]
		}

		if oldorder.Whookupdate == 0 {
			oldorder.Whookupdate = oldorder.Modified
		}
		oldorder.ShipperId = order.ShipperId
		oldorder.Note = order.Note
		oldorder.PartnerShipFee = order.PartnerShipFee
		oldorder.ShipFee = order.ShipFee
		oldorder.ShipmentCode = order.ShipmentCode
		oldorder.Total = order.Total
		oldorder.IsPaid = order.IsPaid
		//oldorder.SearchIndex = mystring.ParameterizeJoin(oldorder.Name+oldorder.Email+oldorder.Phone+oldorder.City+oldorder.District+oldorder.Ward+oldorder.Address+oldorder.CusNote+oldorder.Note+oldorder.ShipmentCode+oldorder.CampaignName+mapshipper[oldorder.ShipperId], " ")
		oldorder.SearchIndex = mystring.ParameterizeJoin(oldorder.Name+oldorder.Email+oldorder.Phone+oldorder.City+oldorder.District+oldorder.Ward+oldorder.Address+oldorder.CusNote+oldorder.Note+oldorder.ShipmentCode+oldorder.CampaignName, " ")
		if mapstat[oldorder.Status].Finish {
			oldorder.Whookupdate = time.Now().Unix()
		}

		rpch.SaveOrder(oldorder)
		oldorder.SearchIndex = ""
		info, _ := json.Marshal(oldorder)
		strrt := string(info)
		return models.RequestResult{Status: 1, Data: strrt}

	}
	return models.RequestResult{Status: 0}

}
func main() {
	//default port for service
	var port string
	port = os.Getenv("PORT")
	if port == "" {
		port = "8903"
	}
	//open service and listen
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Errorf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	fmt.Printf("listening on %s\n", port)
	pb.RegisterGRPCServicesServer(s, &service{})
	if err := s.Serve(lis); err != nil {
		log.Errorf("failed to serve : %v", err)
	}
	fmt.Print("exit")

}
