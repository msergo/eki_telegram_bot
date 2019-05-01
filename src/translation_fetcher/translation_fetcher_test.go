package translation_fetcher

import (
	"testing"
	"github.com/PuerkitoBio/goquery"
	"strings"
)

func TestGetSingleArticle(t *testing.T) {
	rawHtml := "<div class=\"tervikart\"><p> <span id=\"x43640_1_m\" class=\"m x_m m leitud_id\" lang=\"et\"><span class=\"leitud_ss\">probleem</span>+</span><br> <span id=\"x43640_2_x\" class=\"x x_x x\" lang=\"ru\">пробл<span class=\"rohk\">е</span>мный</span>  <span class=\"vormid\">&lt;<span id=\"x43640_3_vgvormid\" class=\"vgvormid x_vgvormid vgvormid\" lang=\"ru\">пробл<span class=\"rohk\">е</span>мная, пробл<span class=\"rohk\">е</span>мное</span>&gt;</span><br> <span id=\"x43640_4_n\" class=\"n x_n n\" lang=\"et\">probleemartikkel</span> <span id=\"x43640_5_qn\" class=\"qn x_qn qn\" lang=\"ru\">пробл<span class=\"rohk\">е</span>мная стать<span class=\"rohk\">я</span></span><br> <span id=\"x43640_6_n\" class=\"n x_n n\" lang=\"et\">probleemkabe</span> <span id=\"x43640_7_v\" class=\"v x_v v\" lang=\"et\">sport</span> <span id=\"x43640_8_qn\" class=\"qn x_qn qn\" lang=\"ru\">композици<span class=\"rohk\">о</span>нные ш<span class=\"rohk\">а</span>шки</span><br> <span id=\"x43640_9_n\" class=\"n x_n n\" lang=\"et\">probleemromaan</span> <span id=\"x43640_10_qn\" class=\"qn x_qn qn\" lang=\"ru\">пробл<span class=\"rohk\">е</span>мный ром<span class=\"rohk\">а</span>н</span><br> <span id=\"x43640_11_n\" class=\"n x_n n\" lang=\"et\">probleemsaade</span> <span id=\"x43640_12_qn\" class=\"qn x_qn qn\" lang=\"ru\">пробл<span class=\"rohk\">е</span>мная перед<span class=\"rohk\">а</span>ча</span><br> <span id=\"x43640_13_n\" class=\"n x_n n\" lang=\"et\">probleemõpe</span> <span id=\"x43640_14_v\" class=\"v x_v v\" lang=\"et\">ped</span> <span id=\"x43640_15_qn\" class=\"qn x_qn qn\" lang=\"ru\">пробл<span class=\"rohk\">е</span>мное обуч<span class=\"rohk\">е</span>ние</span></p></div>"
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(rawHtml))
	got := GetSingleArticle(doc.Nodes[0])
	//expected := "<b>probleem+</b>\r\n" +
	//	"проблемный\r\n" +
	//	"probleemartikkel - проблемная статья\r\n" +
	//	"probleemkabe - композиционные шашки\r\n" +
	//	"probleemromaan - проблемный роман\r\n" +
	//	"probleemsaade - проблемная передача\r\n" +
	//	"probleemõpe - проблемное обучение" //remove exmples for now
	expected := "<b>probleem+</b>\r\n" +
		"проблемный"

	if got != expected {
		t.Errorf("bad card:\r\n %s\r\n%s", got, expected)
	}
}
