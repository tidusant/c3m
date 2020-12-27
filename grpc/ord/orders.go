package main

//repush
import (
	"context"
	"fmt"
	"github.com/tidusant/c3m/common/c3mcommon"
	"github.com/tidusant/c3m/common/mycrypto"
	"github.com/tidusant/c3m/common/mystring"
	maingrpc "github.com/tidusant/c3m/grpc"
	"google.golang.org/grpc"
	"os"
	"time"

	"github.com/tidusant/c3m/common/log"
	"go.mongodb.org/mongo-driver/bson/primitive"

	pb "github.com/tidusant/c3m/grpc/protoc"
	"github.com/tidusant/c3m/repo/cuahang"
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

//extend class MainRPC
type myRPC struct {
	maingrpc.MainRPC
	rpch cuahang.Repo
}

func (s *service) Call(ctx context.Context, in *pb.RPCRequest) (*pb.RPCResponse, error) {
	m := myRPC{}
	//generate user information into usex by calling parent func (m *myRPC) InitUsex that return error string
	rs := models.RequestResult{Error: m.InitUsex(ctx, in, name, ver)}
	//if not error then continue call func
	if rs.Error == "" {
		if m.Usex.Action == "all" {
			rs = m.LoadAll()
		} else if m.Usex.Action == "statusc" {
			rs = m.LoadAllStatusCount()
		} else if m.Usex.Action == "status" {
			rs = m.LoadAllStatus()
		} else if m.Usex.Action == "lao" {
			rs = m.LoadAllOrderByStatus()
		} else if m.Usex.Action == "lg" {
			rs = m.LoadCities()
		} else if m.Usex.Action == "us" {
			rs = m.UpdateOrderStatus()
		} else if m.Usex.Action == "ds" {
			rs = m.DeleteOrderStatus()
		} else if m.Usex.Action == "ss" {
			rs = m.SaveStatus()
		} else if m.Usex.Action == "uo" {
			rs = m.UpdateOrder()
		} else {
			return m.ReturnNilRespone(), nil
		}
	}
	//if there is other repo than rpch then accummulate query count here
	//m.QueryCount+=m.m.Rpch.QueryCount
	return m.ReturnRespone(rs), nil
}

//parse order from web (status="") and update
func (m *myRPC) parseOrder(order models.Order, defaultstatus models.OrderStatus) {
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
		prod := m.Rpch.GetProdByCode(m.Usex.ShopID, code)
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
	m.Rpch.SaveOrder(order)
}
func (m *myRPC) LoadAllOrderByStatus() models.RequestResult {
	log.Debugf("params: %s", m.Usex.Params)
	args := strings.Split(m.Usex.Params, ",")
	statusid, _ := primitive.ObjectIDFromHex(args[0])
	page := 1
	pagesize := 10
	if len(args) > 1 {
		page, _ = strconv.Atoi(args[1])
	}
	if len(args) > 2 {
		size, _ := strconv.Atoi(args[2])
		pagesize = size
	}
	searchterm := ""
	if len(args) > 3 {
		searchterm = args[3]
		if searchterm != "" {
			byteDecode, _ := base64.StdEncoding.DecodeString(mycrypto.Base64fix(searchterm))
			searchterm = string(byteDecode)
			log.Debugf("searchterm: %s", searchterm)
		}
	}

	start := time.Now()
	orders, total, totalPage := m.GetAllOrder(statusid, page, pagesize, searchterm)
	log.Debugf("GetAllOrder: %s", time.Since(start))
	var cus []models.Customer
	if len(orders) > 0 {
		var phones []string
		for _, v := range orders {
			phones = append(phones, v.Phone)
		}
		//map customer
		start = time.Now()
		cus = m.Rpch.GetCustomerByPhones(phones, m.Usex.ShopID)
		log.Debugf("GetCustomerByPhones: %s", time.Since(start))
		start = time.Now()
		cuscount := m.Rpch.CoundOrderByPhones(phones, m.Usex.ShopID)
		log.Debugf("CoundOrderByPhones: %s", time.Since(start))
		log.Debugf("phones count:%d", len(phones))
		log.Debugf("cus count:%d", len(cus))
		log.Debugf("cuscount count:%d", len(cuscount))
		for k, _ := range cus {
			cus[k].OrderCount = cuscount[cus[k].Phone]
		}
	}
	rtdata := struct {
		Total     int
		Orders    []models.Order
		Customers []models.Customer
		PageCount int
	}{Orders: orders, Customers: cus, Total: total, PageCount: totalPage}

	info, err := json.Marshal(rtdata)
	if err != nil {
		fmt.Println(err.Error())
	}
	strrt := string(info)
	return models.RequestResult{Status: 1, Error: "", Data: strrt}
}
func (m *myRPC) LoadAllStatus() models.RequestResult {

	//default status
	status := m.Rpch.GetAllOrderStatus(m.Usex.ShopID)

	info, _ := json.Marshal(status)

	strrt := string(info)
	return models.RequestResult{Status: 1, Error: "", Data: strrt}
}
func (m *myRPC) LoadAllStatusCount() models.RequestResult {

	//default status
	status := m.Rpch.GetAllOrderStatus(m.Usex.ShopID)
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
		status[k].OrderCount = int(m.Rpch.CountOrdersByStatus(m.Usex.ShopID, v.ID, ""))
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
func (m *myRPC) LoadCities() models.RequestResult {
	//
	////default status
	//source:=m.Usex.Params
	//cities := m.Rpch.GetCities(source)
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

func (m *myRPC) UpdateOrderStatus() models.RequestResult {

	info := strings.Split(m.Usex.Params, ",")
	cancelPartner := "0"
	if len(info) > 1 {
		changestatusid, _ := primitive.ObjectIDFromHex(info[1])
		orderid, _ := primitive.ObjectIDFromHex(info[0])

		//check cancel ghtk status:
		status := m.Rpch.GetStatusByID(changestatusid, m.Usex.ShopID)
		ghtkstatussync := status.PartnerStatus["ghtk"]
		if ghtkstatussync != nil {
			for _, statcode := range ghtkstatussync {
				if statcode == "-1" {
					cancelPartner = "1"
				}
			}
		}
		//check stock
		order := m.Rpch.GetOrderByID(orderid, m.Usex.ShopID)
		if order.Status == "" {
			return models.RequestResult{Status: 0, Error: "order not found", Data: ""}
		}
		statusid, _ := primitive.ObjectIDFromHex(order.Status)
		oldstat := m.Rpch.GetStatusByID(statusid, m.Usex.ShopID)
		if oldstat.Export != status.Export {
			//update stock
			sign := 1
			if oldstat.Export {
				sign = -1 //return to stock
			}
			var exportitems []models.ExportItem
			for _, v := range order.Items {
				prod := m.Rpch.GetProdByCode(m.Usex.ShopID, v.ProdCode)
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
						newexportitem.ShopId = m.Usex.ShopID.Hex()
						break
					}
				}
				exportitems = append(exportitems, newexportitem)

			}
			if !m.Rpch.ExportItem(exportitems) {
				return models.RequestResult{Status: 0, Error: "Error update stock", Data: ""}

			}
		}
		m.Rpch.UpdateOrderStatus(m.Usex.ShopID, changestatusid, orderid.Hex())

	}
	return models.RequestResult{Status: 1, Error: "", Message: cancelPartner, Data: ""}

}

func (m *myRPC) SaveStatus() models.RequestResult {

	var status models.OrderStatus
	err := json.Unmarshal([]byte(m.Usex.Params), &status)
	if !c3mcommon.CheckError("update status parse json", err) {
		return models.RequestResult{Status: 0, Error: "update status fail", Data: ""}

	}
	//check old status
	oldstat := status
	if oldstat.ID != primitive.NilObjectID {
		oldstat = m.Rpch.GetStatusByID(status.ID, m.Usex.ShopID)
		oldstat.Title = status.Title
		oldstat.Color = status.Color
		oldstat.Default = status.Default
		oldstat.Finish = status.Finish
		oldstat.Export = status.Export
		oldstat.PartnerStatus = status.PartnerStatus
	} else {
		oldstat.UserId = m.Usex.UserID.Hex()
		oldstat.ShopId = m.Usex.ShopID.Hex()
	}

	//check default
	if oldstat.Default == true {
		m.Rpch.UnSetStatusDefault(m.Usex.ShopID)
	}
	if oldstat.Color == "" {
		oldstat.Color = "ffffff"
	}

	oldstat = m.Rpch.SaveOrderStatus(oldstat)
	b, _ := json.Marshal(oldstat)
	return models.RequestResult{Status: 1, Error: "", Data: string(b)}

}

func (m *myRPC) DeleteOrderStatus() models.RequestResult {
	//get stat
	statusid, _ := primitive.ObjectIDFromHex(m.Usex.Params)
	stat := m.Rpch.GetStatusByID(statusid, m.Usex.ShopID)
	if stat.ID == primitive.NilObjectID {
		return models.RequestResult{Status: 0, Error: "Status not found"}

	}
	if stat.Default {
		return models.RequestResult{Status: 0, Error: "Status is default. Please select another status default."}

	}
	//check status empty
	count := m.Rpch.GetCountOrderByStatus(stat)
	//check old status
	if count > 0 {
		return models.RequestResult{Status: 0, Error: "Status not empty. " + strconv.Itoa(int(count)) + " orders use this status"}

	}

	m.Rpch.DeleteOrderStatus(stat)
	return models.RequestResult{Status: 1, Error: ""}
}

func (m *myRPC) UpdateOrder() models.RequestResult {
	shop := m.Rpch.GetShopById(m.Usex.UserID, m.Usex.ShopID)
	if shop.Status == 0 {
		return models.RequestResult{Status: 0, Error: "Shop is disabled"}

	}
	var order models.Order
	err := json.Unmarshal([]byte(m.Usex.Params), &order)
	if !c3mcommon.CheckError("update order parse json", err) {
		return models.RequestResult{Status: 0, Error: "update order fail"}

	}
	oldorder := order
	mapolditems := make(map[string]models.OrderItem)
	if order.ID != primitive.NilObjectID {
		oldorder = m.Rpch.GetOrderByID(order.ID, shop.ID)
		for _, v := range oldorder.Items {
			mapolditems[v.Code] = v
		}
	} else {
		oldorder.ShopId = m.Usex.ShopID.Hex()
		oldorder.ID = primitive.NewObjectID()
		oldorder.Created = time.Now().Unix()
		oldorder.Modified = oldorder.Created
	}

	//all campaign
	camps := m.Rpch.GetAllCampaigns(m.Usex.ShopID)
	mapcamp := make(map[string]string)
	for _, v := range camps {
		mapcamp[v.ID.Hex()] = v.Name
	}
	//all shipper
	//shippers := m.Rpch.GetAllShipper(m.Usex.ShopID)
	//mapshipper := make(map[string]string)
	//for _, v := range shippers {
	//	mapshipper[v.ID] = v.Name
	//}
	stats := m.Rpch.GetAllOrderStatus(m.Usex.ShopID)
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
			prod := m.Rpch.GetProdByCode(m.Usex.ShopID, v.ProdCode)
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

					newexportitem.ShopId = m.Usex.ShopID.Hex()
					break
				}
			}
			exportitems = append(exportitems, newexportitem)

		}

		//return item stock
		var oldexportitems []models.ExportItem
		for _, v := range mapolditems {
			var old models.ExportItem
			old.ShopId = m.Usex.ShopID.Hex()
			old.Code = v.ProdCode
			old.ItemCode = v.Code
			old.Num = -v.Num
			oldexportitems = append(oldexportitems, old)
		}
		m.Rpch.ExportItem(oldexportitems)
		if !m.Rpch.ExportItem(exportitems) {
			return models.RequestResult{Status: 0, Error: "Error update stock"}

		}

	}
	var cus models.Customer
	if oldorder.Phone == order.Phone {
		cus = m.Rpch.GetCusByPhone(order.Phone, shop.ID.Hex())
	} else if oldorder.Phone == "" && oldorder.Email == order.Email {
		cus = m.Rpch.GetCusByEmail(order.Email, shop.ID.Hex())
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
	if m.Rpch.SaveCus(cus) {
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

		m.Rpch.SaveOrder(oldorder)
		oldorder.SearchIndex = ""
		info, _ := json.Marshal(oldorder)
		strrt := string(info)
		return models.RequestResult{Status: 1, Data: strrt}

	}
	return models.RequestResult{Status: 0}

}

func (m *myRPC) LoadAll() models.RequestResult {
	log.Debugf("params: %s", m.Usex.Params)
	args := strings.Split(m.Usex.Params, ",")
	statusid, _ := primitive.ObjectIDFromHex(args[0])
	page := 1
	pagesize := 10
	if len(args) > 1 {
		page, _ = strconv.Atoi(args[1])
	}
	if len(args) > 2 {
		size, _ := strconv.Atoi(args[2])
		pagesize = size
	}
	searchterm := ""
	if len(args) > 3 {
		searchterm = args[3]
		if searchterm != "" {
			byteDecode, _ := base64.StdEncoding.DecodeString(mycrypto.Base64fix(searchterm))
			searchterm = string(byteDecode)
			log.Debugf("searchterm: %s", searchterm)
		}
	}
	start := time.Now()
	status := m.Rpch.GetAllOrderStatus(m.Usex.ShopID)
	log.Debugf("GetAllOrderStatus: %s", time.Since(start))
	start = time.Now()
	orders, total, totalPage := m.GetAllOrder(statusid, page, pagesize, searchterm)
	log.Debugf("GetAllOrder: %s", time.Since(start))
	var cus []models.Customer
	if len(orders) > 0 {
		var phones []string
		for _, v := range orders {
			phones = append(phones, v.Phone)
		}
		//map customer
		start = time.Now()
		cus = m.Rpch.GetCustomerByPhones(phones, m.Usex.ShopID)
		log.Debugf("GetCustomerByPhones: %s", time.Since(start))
		start = time.Now()
		cuscount := m.Rpch.CoundOrderByPhones(phones, m.Usex.ShopID)
		log.Debugf("CoundOrderByPhones: %s", time.Since(start))
		log.Debugf("phones count:%d", len(phones))
		log.Debugf("cus count:%d", len(cus))
		log.Debugf("cuscount count:%d", len(cuscount))
		for k, _ := range cus {
			cus[k].OrderCount = cuscount[cus[k].Phone]
		}
	}
	rtdata := struct {
		Total     int
		Status    []models.OrderStatus
		Orders    []models.Order
		Customers []models.Customer
		PageCount int
	}{Status: status, Orders: orders, Customers: cus, Total: total, PageCount: totalPage}

	info, err := json.Marshal(rtdata)
	if err != nil {
		fmt.Println(err.Error())
	}
	strrt := string(info)
	return models.RequestResult{Status: 1, Error: "", Data: strrt}

}

func (m *myRPC) GetAllOrder(statusid primitive.ObjectID, page, pagesize int, searchterm string) (orders []models.Order, total int, totalPage int) {

	total = int(m.Rpch.CountOrdersByStatus(m.Usex.ShopID, statusid, searchterm))
	totalPage = (int)(math.Ceil(float64(total) / float64(pagesize)))
	if page > totalPage {
		page = totalPage
	}
	//camps := m.Rpch.GetAllCampaigns(m.Usex.ShopID)
	//mapcamp := make(map[string]string)
	//for _, v := range camps {
	//	mapcamp[v.ID.Hex()] = v.Name
	//}
	////for _, order := range orders {
	////	parseOrder(order, m.Usex, defaultstatus)
	////}

	orders = m.Rpch.GetOrdersByStatus(m.Usex.ShopID, statusid, page, pagesize, searchterm)
	return
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
