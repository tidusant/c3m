package main

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"github.com/tidusant/c3m/common/c3mcommon"
	"github.com/tidusant/c3m/common/mycrypto"
	"github.com/tidusant/c3m/common/mystring"
	maingrpc "github.com/tidusant/c3m/grpc"
	pb "github.com/tidusant/c3m/grpc/protoc"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"image"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"context"

	"google.golang.org/grpc"
	"os"

	"github.com/tidusant/c3m/common/log"

	"github.com/tidusant/c3m/repo/models"

	//	"c3m/common/inflect"
	//	"c3m/log"
	"encoding/json"

	"fmt"
	"net"
)

const (
	name           string = "lptpl"
	ver            string = "1"
	templateFolder        = "./templates"
)

var (
	cdnURL = os.Getenv("CDNURL")
)

type service struct {
	pb.UnimplementedGRPCServicesServer
}

//extend class MainRPC
type myRPC struct {
	maingrpc.MainRPC
}

func (s *service) Call(ctx context.Context, in *pb.RPCRequest) (rt *pb.RPCResponse, err error) {
	m := myRPC{}
	//generate user information into usex by calling parent func (m *myRPC) InitUsex that return error string
	rs := models.RequestResult{Error: m.InitUsex(ctx, in, name, ver)}
	err = nil
	defer func() {
		if err := recover(); err != nil {
			ioutil.WriteFile("templates/"+name+".panic.log", []byte(fmt.Sprint(time.Now().Format("2006-01-02 15:04:05")+" >> panic occurred:", err)), 0644)
			rs.Error = "Something wrong"
			rt = m.ReturnRespone(rs)
		}

	}()
	//if not error then continue call func
	if rs.Error == "" {
		if m.Usex.Action == "s" {
			rs = m.Submit(false)
		} else if m.Usex.Action == "rs" {
			rs = m.Submit(true)
		} else if m.Usex.Action == "rej" {
			rs = m.Reject()
		} else if m.Usex.Action == "ok" {
			rs = m.Approve()
		} else if m.Usex.Action == "lat" {
			rs = m.LoadAllTest()
		} else if m.Usex.Action == "ltpl" {
			rs = m.LoadTemplate()
		} else if m.Usex.Action == "la" {
			rs = m.LoadForBuilder()
		} else if m.Usex.Action == "latpl" {
			rs = m.LoadForUser()
		} else {
			//unknow action
			rt = m.ReturnNilRespone()
			return
		}
	}
	rt = m.ReturnRespone(rs)
	return
}

//load all template for test and approve
func (m *myRPC) LoadAllTest() models.RequestResult {
	if ok, _ := m.Usex.Modules["c3m-lptpl-admin"]; !ok {
		return models.RequestResult{Error: "permission denied"}
	}
	templates, err := m.Rpch.GetAllLpTemplate(m.Usex.UserID, true)
	if err != nil {
		return models.RequestResult{Error: err.Error()}
	}
	var ptemplates []models.LPTemplate
	var atemplates []models.LPTemplate
	for _, v := range templates {
		if v.Status == 2 {
			ptemplates = append(ptemplates, v)
		} else if v.Status == 1 {
			atemplates = append(atemplates, v)
		}
	}
	//b, _ := json.Marshal(ptemplates)
	type RT struct {
		Ptemplates []models.LPTemplate
		Atemplates []models.LPTemplate
	}

	b, err := json.Marshal(RT{Ptemplates: ptemplates, Atemplates: atemplates})
	if err != nil {
		log.Debugf("error:%s", err.Error())
	}
	return models.RequestResult{Status: 1, Data: string(b)}
}

//load all template for builder
func (m *myRPC) LoadForBuilder() models.RequestResult {
	if ok, _ := m.Usex.Modules["c3m-lptpl-builder"]; !ok {
		return models.RequestResult{Error: "permission denied"}
	}
	templates, err := m.Rpch.GetAllLpTemplate(m.Usex.UserID, false)
	if err != nil {
		return models.RequestResult{Error: err.Error()}
	}
	b, _ := json.Marshal(templates)
	return models.RequestResult{Status: 1, Data: string(b)}
}

//load all templates for user
func (m *myRPC) LoadForUser() models.RequestResult {
	if ok, _ := m.Usex.Modules["c3m-lptpl-user"]; !ok {
		return models.RequestResult{Error: "permission denied"}
	}
	templates, err := m.Rpch.GetAllLpTPLForUser()
	if err != nil {
		return models.RequestResult{Error: err.Error()}
	}
	b, _ := json.Marshal(templates)
	return models.RequestResult{Status: 1, Data: string(b), Compress: true}
}
func (m *myRPC) Submit(resubmit bool) models.RequestResult {
	if ok, _ := m.Usex.Modules["c3m-lptpl-builder"]; !ok {
		return models.RequestResult{Error: "permission denied"}
	}
	args := strings.Split(m.Usex.Params, ",")
	if len(args) < 2 {
		return models.RequestResult{Error: "invalid params"}
	}
	tplname := args[0]
	b64str := args[1]
	gzipbyte, err := base64.StdEncoding.DecodeString(b64str)
	if err != nil {
		return models.RequestResult{Error: err.Error()}
	}
	//unzip
	var bb bytes.Buffer
	bb.Write(gzipbyte)

	r, err := gzip.NewReader(&bb)
	if err != nil {
		return models.RequestResult{Error: err.Error()}
	}
	r.Close()
	s, err := ioutil.ReadAll(r)
	if err != nil {
		return models.RequestResult{Error: err.Error()}
	}
	mfile := make(map[string][]byte)
	json.Unmarshal(s, &mfile)
	tplpath := mycrypto.EncodeA(tplname + "_" + m.Usex.Username + mycrypto.StringRand(2))
	os.Mkdir(templateFolder+"/"+tplpath, 0755)
	for k, v := range mfile {
		//check file folder
		if strings.Index(k, "/") > 0 {
			fpath := k[0:strings.LastIndex(k, "/")]
			os.MkdirAll(templateFolder+"/"+tplpath+"/"+fpath, 0755)
		}
		err := ioutil.WriteFile(templateFolder+"/"+tplpath+"/"+k, v, 0644)
		if err != nil {
			return models.RequestResult{Error: err.Error()}
		}
	}

	//update database

	if !resubmit {
		err := m.Rpch.CreateLpTemplate(m.Usex.UserID, tplname, tplpath)
		if err != nil {
			os.RemoveAll(templateFolder + "/" + tplpath)
			return models.RequestResult{Error: err.Error()}
		}
	} else {
		oldtpl, err := m.Rpch.GetLpTemplate(m.Usex.UserID, tplname)
		//reset to waiting approve
		oldtpl.Status = 2
		oldtpl.Path = tplpath
		err = m.Rpch.UpdateLpTemplate(oldtpl)
		if err != nil {
			return models.RequestResult{Error: err.Error()}
		}
	}
	return models.RequestResult{Status: 1, Data: ""}
}

func (m *myRPC) Reject() models.RequestResult {
	if ok, _ := m.Usex.Modules["c3m-lptpl-admin"]; !ok {
		return models.RequestResult{Error: "permission denied"}
	}
	tplID, err := primitive.ObjectIDFromHex(m.Usex.Params)

	if err != nil {
		return models.RequestResult{Error: "template not found"}
	}
	oldtpl, err := m.Rpch.GetLpTemplateById(tplID)
	if err != nil {
		return models.RequestResult{Error: err.Error()}
	}
	//delete content
	os.RemoveAll(templateFolder + "/" + oldtpl.Path)

	//update database
	oldtpl.Status = -1
	oldtpl.Path = ""
	err = m.Rpch.UpdateLpTemplate(oldtpl)
	if err != nil {
		return models.RequestResult{Error: err.Error()}
	}

	return models.RequestResult{Status: 1, Data: ""}
}
func (m *myRPC) Approve() models.RequestResult {

	if ok, _ := m.Usex.Modules["c3m-lptpl-admin"]; !ok {
		return models.RequestResult{Error: "permission denied"}
	}
	tplID, err := primitive.ObjectIDFromHex(m.Usex.Params)

	if err != nil {
		return models.RequestResult{Error: "template not found"}
	}
	oldtpl, err := m.Rpch.GetLpTemplateById(tplID)
	if err != nil {
		return models.RequestResult{Error: err.Error()}
	}

	if oldtpl.Status != 2 {
		return models.RequestResult{Error: "template not in review status"}
	}

	//=======================zip file
	//tmplFolder := templateFolder + `/` + oldtpl.Path
	//zipfile, err := os.Create(templateFolder + `/` + oldtpl.Path + `.zip`)
	//if err != nil {
	//	return models.RequestResult{Error: err.Error()}
	//}
	//defer zipfile.Close()
	//
	//archive := zip.NewWriter(zipfile)
	//defer archive.Close()
	//
	//info, err := os.Stat(tmplFolder)
	//if err != nil {
	//	return models.RequestResult{Error: err.Error()}
	//}
	//
	//var baseDir string
	//if info.IsDir() {
	//	baseDir = filepath.Base(tmplFolder)
	//}
	//
	//filepath.Walk(tmplFolder, func(path string, info os.FileInfo, err error) error {
	//	if err != nil {
	//		return err
	//	}
	//	//skip folder images
	//
	//	if path == tmplFolder+"/content.html" || path == tmplFolder+"/items.html" || path == tmplFolder+"/navitem.html" || strings.Index(path, tmplFolder+"/css") == 0 || strings.Index(path, tmplFolder+"/js") == 0 || strings.Index(path, tmplFolder+"/itemicons") == 0 {
	//
	//		header, err := zip.FileInfoHeader(info)
	//		if err != nil {
	//			return err
	//		}
	//
	//		if baseDir != "" {
	//			header.Name = filepath.Join(baseDir, strings.TrimPrefix(path, tmplFolder))
	//		}
	//
	//		if info.IsDir() {
	//			header.Name += "/"
	//		} else {
	//			header.Method = zip.Deflate
	//		}
	//
	//		writer, err := archive.CreateHeader(header)
	//		if err != nil {
	//			return err
	//		}
	//
	//		if info.IsDir() {
	//			return nil
	//		}
	//
	//		file, err := os.Open(path)
	//		if err != nil {
	//			return err
	//		}
	//		defer file.Close()
	//		_, err = io.Copy(writer, file)
	//		return err
	//	}
	//	return nil
	//})
	//if err != nil {
	//	return models.RequestResult{Error: err.Error()}
	//}
	//=======================

	//======publish screen shot
	tmplFolder := templateFolder + `/` + oldtpl.Path
	lsCDNImages := []string{}
	cdnfolder := "./cdn/lptemplate"
	os.MkdirAll(cdnfolder, 0755)
	filename := mycrypto.StringRand(5) + mycrypto.StringRand(5) + mycrypto.StringRand(5) + ".jpg"

	//read file
	file, err := os.Open(tmplFolder + "/screenshot.jpg")
	if err != nil {
		file.Close()
		return models.RequestResult{Error: "error reading screenshot"}
	}
	imageconfig, _, _ := image.DecodeConfig(file)
	file.Close()
	fileb, err := ioutil.ReadFile(tmplFolder + "/screenshot.jpg")
	if err != nil {
		return models.RequestResult{Error: "error reading screenshot"}
	}
	//copy file
	err = ioutil.WriteFile(cdnfolder+"/"+filename, fileb, 0644)
	if err != nil {
		return models.RequestResult{Error: "error creating screenshot"}
	}
	lsCDNImages = append(lsCDNImages, filename)
	//create thumb
	thumbwidth := 200
	if imageconfig.Width > thumbwidth {
		fileb, _ = c3mcommon.ImgResize(fileb, uint(thumbwidth), 0)
		err = ioutil.WriteFile(cdnfolder+"/thumb_"+filename, fileb, 0644)
		if err != nil {
			m.removeCDNImages(lsCDNImages)
			return models.RequestResult{Error: "error creating thumb screenshot"}
		}
	}
	lsCDNImages = append(lsCDNImages, "thumb_"+filename)

	//copy all image into cdn
	files, _ := ioutil.ReadDir(tmplFolder)
	for _, f := range files {
		if !f.IsDir() {
			fname := f.Name()
			if fname == "content.html" || fname == "items.html" || fname == "navitem.html" {
				images, errstr := m.replaceCDNImages(tmplFolder+"/"+fname, cdnfolder, `<img.*src="(.*?)".*>`, `{{templatePath}}`, tmplFolder)
				lsCDNImages = append(lsCDNImages, images...)
				if errstr != "" {
					m.removeCDNImages(lsCDNImages)
					return models.RequestResult{Error: errstr}
				}
			}
		}
	}

	files, _ = ioutil.ReadDir(tmplFolder + "/css")
	for _, f := range files {
		if !f.IsDir() {
			fname := f.Name()
			if filepath.Ext(fname) == ".css" {
				images, errstr := m.replaceCDNImages(tmplFolder+"/css/"+fname, cdnfolder, `url\(["']?(.*?)["']?\)`, `../`, tmplFolder+`/`)
				lsCDNImages = append(lsCDNImages, images...)
				if errstr != "" {
					m.removeCDNImages(lsCDNImages)
					return models.RequestResult{Error: errstr}
				}
			}
		}
	}

	log.Debugf("lsCDNImages:%+v", lsCDNImages)

	//=======================
	//update database
	oldtpl.ScreenShot = filename
	oldtpl.Status = 1
	err = m.Rpch.UpdateLpTemplate(oldtpl)
	if err != nil {
		return models.RequestResult{Error: err.Error()}
	}

	return models.RequestResult{Status: 1, Data: ""}
}

func (m *myRPC) replaceCDNImages(filepath, cdnfolder, regex, oldstr, newstr string) (images []string, errstr string) {
	defer func() {
		if err := recover(); err != nil {
			log.Debugf("replaceCDNImages panic: %s", err)
		}
	}()
	errstr = ""
	images = []string{}
	htmlContent := ""
	cdnImageURL := cdnURL + "lptemplate/"
	tplContent, err := ioutil.ReadFile(filepath)
	if err != nil {
		return images, "error reading html content"
	}
	htmlContent = string(tplContent)
	//find, create and replace image in html
	var reg = regexp.MustCompile(regex)
	t := reg.FindAllStringSubmatch(htmlContent, -1)
	for _, v := range t {
		if strings.Index(v[1], "http://") == 0 || strings.Index(v[1], "https://") == 0 {
			continue
		}
		imagePath := strings.Replace(v[1], oldstr, newstr, 1)
		log.Debugf("imagePath %s - %s - %s", imagePath, filepath, v[0])
		if strings.Index(imagePath, "data:image/") > -1 {
			continue
		}
		fileext := imagePath[strings.LastIndex(imagePath, "."):]
		fileb, err := ioutil.ReadFile(imagePath)
		if err != nil {
			return images, "error reading image " + imagePath
		}
		//write image
		filename := mystring.RandString(5) + mystring.RandString(5) + mystring.RandString(5) + fileext
		err = ioutil.WriteFile(cdnfolder+"/"+filename, fileb, 0644)
		if err != nil {
			return images, "error creating image " + imagePath
		}
		images = append(images, filename)
		htmlContent = strings.Replace(htmlContent, v[1], cdnImageURL+filename, 1)
	}
	//write to html file
	err = ioutil.WriteFile(filepath, []byte(htmlContent), 0644)
	if err != nil {
		return images, "error update html file " + filepath
	}
	return
}

func (m *myRPC) removeCDNImages(images []string) {
	cdnfolder := "./cdn/lptemplate"
	for _, v := range images {
		os.Remove(cdnfolder + "/" + v)
	}
}
func (m *myRPC) LoadTemplate() models.RequestResult {
	if ok, _ := m.Usex.Modules["c3m-lptpl-user"]; !ok {
		return models.RequestResult{Error: "permission denied"}
	}
	tplID, err := primitive.ObjectIDFromHex(m.Usex.Params)

	if err != nil {
		return models.RequestResult{Error: "template not found"}
	}
	tpl, err := m.Rpch.GetLpTemplateById(tplID)
	if err != nil {
		return models.RequestResult{Error: err.Error()}
	}
	if tpl.Status != 1 {
		return models.RequestResult{Error: "template not found"}
	}

	mfile := make(map[string][]byte)
	tmplFolder := templateFolder + `/` + tpl.Path
	//check folder exist
	if _, err := os.Stat(tmplFolder); os.IsNotExist(err) {
		return models.RequestResult{Error: "template directory not found"}
	}
	walker := func(path string, info os.FileInfo, err error) error {
		readPath := "templates/" + tpl.Path
		//skip tailwind.css
		if path == readPath+"/css/tailwind.css" {
			return nil
		}
		//skip folder images

		if path == readPath+"/content.html" || path == readPath+"/items.html" || path == readPath+"/navitem.html" || strings.Index(path, readPath+"/css") == 0 || strings.Index(path, readPath+"/js") == 0 || strings.Index(path, readPath+"/itemicons") == 0 {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			b, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			mfile[strings.Replace(path, "templates/"+tpl.Path+"/", "", 1)] = b

		}
		return nil
	}
	err = filepath.Walk(templateFolder+"/"+tpl.Path+"/", walker)
	if err != nil {
		return models.RequestResult{Error: err.Error()}
	}

	b, _ := json.Marshal(mfile)
	if err != nil {
		return models.RequestResult{Error: err.Error()}
	}

	//b2, err := ioutil.ReadFile(templateFolder + "/" + tpl.Path + ".zip")
	//b := base64.StdEncoding.EncodeToString(b2)

	return models.RequestResult{Status: 1, Data: string(b), Compress: true}
}

func main() {
	//default port for service
	var port string
	port = os.Getenv("PORT")
	if port == "" {
		port = "8905"
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
