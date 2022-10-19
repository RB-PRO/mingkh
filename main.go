package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/gosuri/uilive"
	"github.com/xuri/excelize/v2"
)

type region struct {
	link   string
	region string
}
type CityLink struct {
	link   string
	city   string
	region string
}
type zhekLink struct {
	link   string
	name   string
	city   string
	region string
}

type zhek struct {
	link    string // Ссылка на дело
	city    string // Город
	region  string // Регион
	name    string // Название компании
	nameAll string // ПолноеНазвание компании
	inn     string // ИНН
	ogrn    string // ОГРН
	adresUR string // Адрес регистрации ЮЛ
	adresF  string // Фактическое местонахождение органов управления
	tel1    string // Контактные телефоны
	tel2    string // Телефон диспетчерской службы
	email   string // e-mail
	site    string // Официальный сайт
	people  string // ФИО руководителя
	dom     string // Дома в управлении
	metr    string // Метров квадратных
}

const site string = "https://mingkh.ru"

var ind int = 2

func main() {
	fmt.Println("-> Нажмите на Enter, чтобы загрузить список регионов")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	fmt.Println("Загружаю список регионов")

	f_excel := excelize.NewFile()
	makeHeadXlsx(f_excel)

	regionslinks := regionMakingLinks()
	for ind, val := range regionslinks {
		fmt.Printf("%v. %s\n", ind+1, val.region)
	}
	fmt.Printf("-> Выберите номера регионов(1 2 3), которые необходимо загрузить\n")
	fmt.Printf("-> или выберите номер региона(4), в котором необходимо загрузить данные по одноме региону\n")
	fmt.Printf("-> или выберите Нуль(0), чтобы собрать данные по всем регионам.\n")
	myscanner := bufio.NewScanner(os.Stdin)
	myscanner.Scan()
	line := myscanner.Text()
	strs := strings.Split(line, " ")
	var regionslinksNEW []region
	if len(strs) == 0 || strs[0] == "0" {
		fmt.Println("Начинаю парсить все регионы(это будет долго).")
		citylinks := cityMakingLinks(regionslinks)
		zhekLinks := zhekLinkMake(citylinks)
		zhekMakinkList(f_excel, zhekLinks)
	} else if len(strs) > 1 && strs[0] != "0" {
		for i := 0; i < len(strs); i++ {
			var strINT int
			strINT, _ = strconv.Atoi(strs[i])
			if strINT > 0 && strINT <= len(regionslinks) {
				regionslinksNEW = append(regionslinksNEW, region{link: regionslinks[strINT-1].link, region: regionslinks[strINT-1].region})
			}
		}
		fmt.Println("Начинаю парсить все ЖКХ по регионам:")
		for ind, val := range regionslinksNEW {
			fmt.Printf("%v. %s; ", ind+1, val.region)
		}
		citylinks := cityMakingLinks(regionslinksNEW)
		zhekLinks := zhekLinkMake(citylinks)
		zhekMakinkList(f_excel, zhekLinks)

	} else if len(strs) == 1 {
		var strINT int
		strINT, _ = strconv.Atoi(strs[0])
		regionslinksNEW = append(regionslinksNEW, region{link: regionslinks[strINT-1].link, region: regionslinks[strINT-1].region})
		citylinks := cityMakingLinks(regionslinksNEW)
		//fmt.Println(citylinks)
		for ind, val := range citylinks {
			fmt.Printf("%v. %s\n", ind+1, val.city)
		}
		fmt.Printf("-> Выберите номера города(1 2 3), которые необходимо загрузить\n")
		fmt.Printf("-> или выберите Нуль(0), чтобы собрать данные по всем городам.\n")
		myscanner := bufio.NewScanner(os.Stdin)
		myscanner.Scan()
		line := myscanner.Text()
		strs := strings.Split(line, " ")
		//fmt.Println(strs)
		var citylinksNEW []CityLink
		for i := 0; i < len(strs); i++ {
			var strINT int
			strINT, _ = strconv.Atoi(strs[i])
			if strINT > 0 && strINT <= len(citylinks) {
				citylinksNEW = append(citylinksNEW, CityLink{link: citylinks[strINT-1].link, city: citylinks[strINT-1].city})
			}
		}

		if len(strs) == 0 || strs[0] == "0" {
			fmt.Println("Начинаю парсить все ЖКХ по всем городам:")
			for ind, val := range citylinks {
				fmt.Printf("%v. %s; ", ind+1, val.city)
			}
			zhekLinks := zhekLinkMake(citylinks)
			zhekMakinkList(f_excel, zhekLinks)
		} else {
			fmt.Println("Начинаю парсить все ЖКХ по городам:")
			for ind, val := range citylinksNEW {
				fmt.Printf("%v. %s; ", ind+1, val.city)
			}
			zhekLinks := zhekLinkMake(citylinksNEW)
			zhekMakinkList(f_excel, zhekLinks)
		}
	}

	if err := f_excel.SaveAs("zheks.xlsx"); err != nil {
		fmt.Println(err)
	}

}

func zhekMakinkList(f *excelize.File, zhekLinks []zhekLink) {
	var tecZhek, cleartecZhek zhek
	c := colly.NewCollector()
	c.OnHTML("div[class=breads] ul", func(e *colly.HTMLElement) {
		//fmt.Println(1)
		tecZhek.region = e.DOM.Find("li:nth-child(2)").Text()
		tecZhek.city = e.DOM.Find("li:nth-child(3)").Text()
		tecZhek.name = e.DOM.Find("li:nth-child(4)").Text()
		tecZhek.link, _ = e.DOM.Find("li:nth-child(4) a").Attr("href")
		tecZhek.link = site + tecZhek.link
	})
	c.OnHTML("div[id=registracionnye-dannye] div[class^=table-responsive] table[class^=table] tbody tr td", func(e *colly.HTMLElement) {
		//fmt.Println(4)
		if strings.Contains(e.DOM.Text(), "Полное наименование") {
			tecZhek.nameAll = e.DOM.Next().Text()
		}
		if strings.Contains(e.DOM.Text(), "ИНН") {
			tecZhek.inn = e.DOM.Next().Text()
		}
		if strings.Contains(e.DOM.Text(), "ОГРН") {
			tecZhek.ogrn = e.DOM.Next().Text()
		}
		if strings.Contains(e.DOM.Text(), "Адрес регистрации ЮЛ") {
			tecZhek.adresUR = e.DOM.Next().Text()
		}
		//***
		if strings.Contains(e.DOM.Text(), "Фактическое местонахождение органов управления") {
			tecZhek.adresF = e.DOM.Next().Text()
		}
		if strings.Contains(e.DOM.Text(), "Контактные телефоны") {
			tecZhek.tel1 = e.DOM.Next().Text()
		}
		if strings.Contains(e.DOM.Text(), "Телефон диспетчерской службы") {
			tecZhek.tel2 = e.DOM.Next().Text()
		}
		if strings.Contains(e.DOM.Text(), "ФИО руководителя") {
			tecZhek.people = e.DOM.Next().Text()
		}
		if strings.Contains(e.DOM.Text(), "Официальный сайт") {
			tecZhek.site = e.DOM.Next().Text()
		}
	})

	c.OnHTML("div[class=col-md-7]", func(e *colly.HTMLElement) {
		tecZhek.dom = e.DOM.Find("span:nth-child(1)").Text()
		tecZhek.dom = strings.Replace(tecZhek.dom, "&nbsp;", "", -1)
		tecZhek.dom = strings.Replace(tecZhek.dom, " ", "", -1)
		tecZhek.dom = strings.Replace(tecZhek.dom, "дома", "", -1)
		tecZhek.dom = strings.Replace(tecZhek.dom, "домов", "", -1)
		tecZhek.dom = strings.Replace(tecZhek.dom, "дом", "", -1)

		tecZhek.metr = e.DOM.Find("span:nth-child(2)").Text()
		tecZhek.metr = strings.Replace(tecZhek.metr, "&nbsp;", "", -1)
		tecZhek.metr = strings.Replace(tecZhek.metr, " ", "", -1)
		tecZhek.metr = strings.Replace(tecZhek.metr, "м2", "", -1)
	})

	c.OnHTML("dl[class^=dl-horizontal] dt", func(e *colly.HTMLElement) {
		if e.DOM.Text() == "E-mail" {
			tecZhek.email = e.DOM.Next().Text()
		}
		if e.DOM.Text() == "Веб-сайт" {
			tecZhek.site = e.DOM.Next().Text()
		}
	})

	fmt.Println()
	writer := uilive.New()
	// start listening for updates and render
	writer.Start()
	for ind, tea := range zhekLinks {
		c.Visit(site + tea.link)
		writeInXlsx(f, tecZhek)
		tecZhek = cleartecZhek
		fmt.Fprintf(writer, "   Загрузка (%d/%d)\n", ind+1, len(zhekLinks))
	}
	writer.Stop() // flush and stop rendering
}

func zhekLinkMake(citylinks []CityLink) []zhekLink {
	var zhekLinks []zhekLink
	var zheksLinks zhekLink
	c := colly.NewCollector()
	c.OnHTML("tbody a", func(e *colly.HTMLElement) {
		zheksLinks.link, _ = e.DOM.Attr("href")
		//zheksLinks.region = e.DOM.Text()
		zhekLinks = append(zhekLinks, zheksLinks)
		//writeInXlsx(f_excel, tecZhek)
	})
	for _, val := range citylinks {
		c.Visit(site + "/rating" + val.link)
	}
	//fmt.Println(site + "/rating" + citylinks[0].link)
	//c.Visit(site + "/rating" + citylinks[0].link)
	return zhekLinks
}

func regionMakingLinks() []region {
	var regions []region
	var reg region
	c := colly.NewCollector()
	c.OnHTML("ul[class^=col-md-3] li", func(e *colly.HTMLElement) {

		reg.link, _ = e.DOM.Find("a").Attr("href")
		reg.region = e.DOM.Find("a").Text()
		regions = append(regions, reg)
	})
	c.Visit(site)
	return regions
}

func cityMakingLinks(regions []region) []CityLink {

	var citylinks []CityLink
	var city CityLink

	b := colly.NewCollector()
	b.OnHTML("ul[class^=col-md-3] li", func(e *colly.HTMLElement) {
		city.link, _ = e.DOM.Find("a").Attr("href")
		city.city = e.DOM.Find("a").Text()
		citylinks = append(citylinks, city)
	})
	for _, val := range regions {
		b.Visit(site + val.link)
	}
	//b.Visit(site + regions[0].link)
	return citylinks
}
func writeInXlsxArr(f *excelize.File, vals []zhek) {
	for _, val := range vals {
		writeInXlsx(f, val)
	}
}
func writeInXlsx(f *excelize.File, val zhek) {

	f.SetCellValue("main", "A"+strconv.Itoa(ind), val.link)
	f.SetCellValue("main", "B"+strconv.Itoa(ind), val.city)
	f.SetCellValue("main", "C"+strconv.Itoa(ind), val.region)
	f.SetCellValue("main", "D"+strconv.Itoa(ind), val.name)
	f.SetCellValue("main", "E"+strconv.Itoa(ind), val.inn)
	f.SetCellValue("main", "F"+strconv.Itoa(ind), val.ogrn)
	f.SetCellValue("main", "G"+strconv.Itoa(ind), val.adresUR)
	f.SetCellValue("main", "H"+strconv.Itoa(ind), val.adresF)
	f.SetCellValue("main", "I"+strconv.Itoa(ind), val.tel1)
	f.SetCellValue("main", "J"+strconv.Itoa(ind), val.tel2)
	f.SetCellValue("main", "K"+strconv.Itoa(ind), val.email)
	f.SetCellValue("main", "L"+strconv.Itoa(ind), val.site)
	f.SetCellValue("main", "M"+strconv.Itoa(ind), val.people)
	f.SetCellValue("main", "N"+strconv.Itoa(ind), val.dom)
	f.SetCellValue("main", "O"+strconv.Itoa(ind), val.metr)
	ind++
}

func makeHeadXlsx(f *excelize.File) {
	f.NewSheet("main")
	f.DeleteSheet("Sheet1")
	f.SetCellValue("main", "A1", "Ссылка на дело")
	f.SetCellValue("main", "B1", "Город")
	f.SetCellValue("main", "C1", "Регион")
	f.SetCellValue("main", "D1", "Название компании")
	f.SetCellValue("main", "E1", "ИНН")
	f.SetCellValue("main", "F1", "ОГРН")
	f.SetCellValue("main", "G1", "Адрес регистрации ЮЛ")
	f.SetCellValue("main", "H1", "Фактическое местонахождение органов управления")
	f.SetCellValue("main", "I1", "Контактные телефоны")
	f.SetCellValue("main", "J1", "Телефон диспетчерской службы")
	f.SetCellValue("main", "K1", "E-mail")
	f.SetCellValue("main", "L1", "Официальный сайт")
	f.SetCellValue("main", "M1", "ФИО руководителя")
	f.SetCellValue("main", "N1", "Дома в управлении")
	f.SetCellValue("main", "O1", "Метров квадратных")
}
