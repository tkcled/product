package service_product

import (
	"context"
	"fmt"
	"hshelby-tkcled-product/src/database/collection"
	model "hshelby-tkcled-product/src/database/model"
	"log"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Crawl() {
	ctx := context.Background()
	categories := []model.Category{}

	cur, err := collection.Category().Collection().Find(ctx, bson.M{})
	if err != nil {
		log.Println("ImportProduct get parentCategories", err)
	}

	err = cur.All(ctx, &categories)
	if err != nil {
		log.Println("ImportProduct cur.All(ctx, &parentCategories)", err)
	}

	currentCategories := map[string][]string{}
	for _, ele := range categories {
		currentCategories[strings.ToLower(ele.Name)] = []string{ele.ID, ele.ParentID}
	}

	links := []string{"https://www.tkcme.com/c/bo-den-hibay-led-philipsc.htm", "https://www.tkcme.com/c/bo-den-pha-led-philipsc.htm", "https://www.tkcme.com/c/bo-den-metal-halide.htm", "https://www.tkcme.com/c/bo-den-sodium-natri.htm", "https://www.tkcme.com/c/bo-den-duong-led-philipsc.htm", "https://www.tkcme.com/c/Bo-den-sodium-natric.htm", "https://www.tkcme.com/c/bo-den-metal.htm", "https://www.tkcme.com/c/bo-den-nlmtc.htm", "https://www.tkcme.com/c/bo-den-panel-philips-certaflux-led-panel-5959c.htm", "https://www.tkcme.com/c/bo-den-panel-philips-certaflux-led-panel-6060c.htm", "https://www.tkcme.com/c/bo-den-panel-philips-certaflux-led-panel-6060-865-840-830-md2c.htm", "https://www.tkcme.com/c/dèn-chóng-thám-certaflux-waterproof.htm", "https://www.tkcme.com/c/dèn-chóng-thám-fortimo-waterproof.htm", "https://www.tkcme.com/c/cao-ap-metalhalidec.htm", "https://www.tkcme.com/c/bong-cdmc.htm", "https://www.tkcme.com/c/ac.htm", "https://www.tkcme.com/c/bong-den-ledtube-t5c.htm", "https://www.tkcme.com/c/bong-den-led-parc.htm", "https://www.tkcme.com/c/chan-luu-ballast-dien-tu-philips.htm", "https://www.tkcme.com/c/ballast-kich-tu-osramc.htm", "https://www.tkcme.com/c/fortimo-slm.htm", "https://www.tkcme.com/c/bong-tuyp-t5c.htm", "https://www.tkcme.com/c/bong-vongc.htm", "https://www.tkcme.com/c/tac-te-con-chuotc.htm", "https://www.tkcme.com/c/chan-luu-dien-tu-t8c.htm", "https://www.tkcme.com/c/chan-luu-dien-tu-t5c.htm", "https://www.tkcme.com/c/uv-c-pond-water.htm", "https://www.tkcme.com/c/dulux-dc.htm", "https://www.tkcme.com/c/dulux-lc.htm", "https://www.tkcme.com/c/sieu-ben.htm", "https://www.tkcme.com/c/loai-co-led-tich-hop.htm", "https://www.tkcme.com/c/loai-dung-voi-chan-luu-dien-tu.htm", "https://www.tkcme.com/c/choa-downlight.htm", "https://www.tkcme.com/c/transformer---bien-the.htm", "https://www.tkcme.com/c/halogen-classic.htm", "https://www.tkcme.com/c/haloline.htm", "https://www.tkcme.com/c/haloline-eco.htm", "https://www.tkcme.com/c/halolux.htm", "https://www.tkcme.com/c/halopar.htm", "https://www.tkcme.com/c/halosport.htm", "https://www.tkcme.com/c/halostar.htm", "https://www.tkcme.com/c/halopin.htm", "https://www.tkcme.com/c/bong-den-essential-ledspot-mr16.htm", "https://www.tkcme.com/c/den-led-bup.htm", "https://www.tkcme.com/c/bong-chanh.htm", "https://www.tkcme.com/c/bong-phan-xa.htm", "https://www.tkcme.com/c/bong-par.htm", "https://www.tkcme.com/c/bong-den-be-ca.htm", "https://www.tkcme.com/c/bong-den-khu-trung.htm", "https://www.tkcme.com/c/bong-den-hong-ngoai.htm", "https://www.tkcme.com/c/bong-san-khau-dien-anh-y-te.htm", "https://www.tkcme.com/c/bong-ho-boi.htm", "https://www.tkcme.com/c/bong-tu-ao---linestra.htm"}

	c := colly.NewCollector()

	mapProduct := map[string]string{}

	for _, link := range links {
		findCategory := false
		categoryID := ""

		if !findCategory {
			c.OnHTML("div.maincol h1", func(e *colly.HTMLElement) {
				category := strings.ToLower(e.Text)
				for _, tmp := range currentCategories[category] {
					if tmp != "" {
						categoryID = fmt.Sprintf("%s, %s", categoryID, tmp)
					}
				}
				findCategory = true
			})

			c.OnHTML("div.productitem", func(e *colly.HTMLElement) {
				name := e.ChildAttr("img", "alt")
				image := e.ChildAttr("img", "src")

				_, exist := mapProduct[name]

				if exist {
					condition := make(map[string]interface{})
					condition["name"] = name

					result := &model.Product{}
					err = collection.Product().Collection().FindOne(ctx, condition).Decode(result)
					if err != nil {
						fmt.Println(name)
						fmt.Println(err)
						return
					}

					updated := make(map[string]interface{})

					updated["category_id"] = fmt.Sprintf("%s,%s", result.CategoryID, categoryID)

					_, err = collection.Product().Collection().UpdateByID(ctx, result.ID, bson.M{"$set": updated})
					if err != nil {
						fmt.Println(err)
						return
					}
				}

				if !exist {
					mapProduct[name] = name
					product := model.Product{
						ID: primitive.NewObjectID().Hex(),

						Name:        name,
						Image:       image,
						Description: name,
						Code:        "",
						CatalogLink: "",
						CategoryID:  categoryID,

						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					}
					_, err = collection.Product().Collection().InsertOne(ctx, product)
					if err != nil {
						fmt.Println(err)
						return
					}

				}

			})

			// Bắt đầu truy cập trang web
			err := c.Visit(link)
			if err != nil {
				fmt.Printf(link)
				log.Fatal(err, link)
			}
		}
	}
	fmt.Println("done")
}
