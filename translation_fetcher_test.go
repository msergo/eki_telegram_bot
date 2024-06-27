package main

import (
	"strings"
	"testing"

	"os"

	"github.com/PuerkitoBio/goquery"
)

func TestGetSingleArticledWithoutGrammarForms(t *testing.T) {
	rawHTML := "<div class=\"tervikart\"><p> <span id=\"x43640_1_m\" class=\"m x_m m leitud_id\" lang=\"et\"><span class=\"leitud_ss\">probleem</span>+</span><br> <span id=\"x43640_2_x\" class=\"x x_x x\" lang=\"ru\">пробл<span class=\"rohk\">е</span>мный</span>  <span class=\"vormid\">&lt;<span id=\"x43640_3_vgvormid\" class=\"vgvormid x_vgvormid vgvormid\" lang=\"ru\">пробл<span class=\"rohk\">е</span>мная, пробл<span class=\"rohk\">е</span>мное</span>&gt;</span><br> <span id=\"x43640_4_n\" class=\"n x_n n\" lang=\"et\">probleemartikkel</span> <span id=\"x43640_5_qn\" class=\"qn x_qn qn\" lang=\"ru\">пробл<span class=\"rohk\">е</span>мная стать<span class=\"rohk\">я</span></span><br> <span id=\"x43640_6_n\" class=\"n x_n n\" lang=\"et\">probleemkabe</span> <span id=\"x43640_7_v\" class=\"v x_v v\" lang=\"et\">sport</span> <span id=\"x43640_8_qn\" class=\"qn x_qn qn\" lang=\"ru\">композици<span class=\"rohk\">о</span>нные ш<span class=\"rohk\">а</span>шки</span><br> <span id=\"x43640_9_n\" class=\"n x_n n\" lang=\"et\">probleemromaan</span> <span id=\"x43640_10_qn\" class=\"qn x_qn qn\" lang=\"ru\">пробл<span class=\"rohk\">е</span>мный ром<span class=\"rohk\">а</span>н</span><br> <span id=\"x43640_11_n\" class=\"n x_n n\" lang=\"et\">probleemsaade</span> <span id=\"x43640_12_qn\" class=\"qn x_qn qn\" lang=\"ru\">пробл<span class=\"rohk\">е</span>мная перед<span class=\"rohk\">а</span>ча</span><br> <span id=\"x43640_13_n\" class=\"n x_n n\" lang=\"et\">probleemõpe</span> <span id=\"x43640_14_v\" class=\"v x_v v\" lang=\"et\">ped</span> <span id=\"x43640_15_qn\" class=\"qn x_qn qn\" lang=\"ru\">пробл<span class=\"rohk\">е</span>мное обуч<span class=\"rohk\">е</span>ние</span></p></div>"
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(rawHTML))
	got, isMainArticle := GetSingleArticleWithDirection("probleem", doc.Nodes[0], "est-rus")
	expected := "<b>probleem+</b>\r\n" +
		"проблемный"

	if got != expected {
		t.Errorf("bad card:\r\n %s\r\n%s", got, expected)
	}
	if isMainArticle != false {
		t.Errorf("article without grammar forms should not be the main")

	}
}

func TestGetSingleArticledWithGrammarForms(t *testing.T) {
	rawHTML := "<div class=\"tervikart\"><p> <span id=\"x43638_1_m\" class=\"m x_m m leitud_id\" lang=\"et\"><span class=\"leitud_ss\">probleem</span></span> <span id=\"x43638_2_sl\" class=\"sl x_sl sl\" lang=\"et\">s</span> <span class=\"vormid\">&lt;<span id=\"x43638_3_mv\" class=\"mv x_mv mv\" lang=\"et\">probl'eem probleemi probl'eemi probl'eemi, probl'eemi[de probl'eemi[sid ~ probl'eem/e</span> <span id=\"x43638_4_mt\" class=\"mt x_mt mt\" lang=\"et\">22</span>&gt;</span><br> <span id=\"x43638_5_d\" class=\"d x_d d\" lang=\"et\">uurimisülesanne</span>; <span id=\"x43638_6_d\" class=\"d x_d d\" lang=\"et\">lahendust nõudev keerukas küsimus</span><br> <span id=\"x43638_7_x\" class=\"x x_x x\" lang=\"ru\">пробл<span class=\"rohk\">е</span>ма</span> <span class=\"vormid\">&lt;<span id=\"x43638_8_vgvormid\" class=\"vgvormid x_vgvormid vgvormid\" lang=\"ru\">пробл<span class=\"rohk\">е</span>мы</span> <span id=\"x43638_9_vgsugu\" class=\"vgsugu x_vgsugu vgsugu\" lang=\"ru\">ж</span>&gt;</span>,<br> <span id=\"x43638_10_x\" class=\"x x_x x\" lang=\"ru\">вопр<span class=\"rohk\">о</span>с</span> <span class=\"vormid\">&lt;<span id=\"x43638_11_vgvormid\" class=\"vgvormid x_vgvormid vgvormid\" lang=\"ru\">вопр<span class=\"rohk\">о</span>са</span> <span id=\"x43638_12_vgsugu\" class=\"vgsugu x_vgsugu vgsugu\" lang=\"ru\">м</span>&gt;</span><br> <span id=\"x43638_13_n\" class=\"n x_n n\" lang=\"et\">lahendamatu probleem</span> <span id=\"x43638_14_qn\" class=\"qn x_qn qn\" lang=\"ru\">неразреш<span class=\"rohk\">и</span>мая пробл<span class=\"rohk\">е</span>ма</span><br> <span id=\"x43638_15_n\" class=\"n x_n n\" lang=\"et\">filosoofiaprobleem</span> <span id=\"x43638_16_qn\" class=\"qn x_qn qn\" lang=\"ru\">филос<span class=\"rohk\">о</span>фская пробл<span class=\"rohk\">е</span>ма</span> /  <span id=\"x43638_17_qn\" class=\"qn x_qn qn\" lang=\"ru\">пробл<span class=\"rohk\">е</span>ма филос<span class=\"rohk\">о</span>фии</span><br> <span id=\"x43638_18_n\" class=\"n x_n n\" lang=\"et\">seksuaalprobleem</span> <span id=\"x43638_19_qn\" class=\"qn x_qn qn\" lang=\"ru\">сексу<span class=\"rohk\">а</span>льная пробл<span class=\"rohk\">е</span>ма</span><br> <span id=\"x43638_20_n\" class=\"n x_n n\" lang=\"et\">teadusprobleem</span> <span id=\"x43638_21_qn\" class=\"qn x_qn qn\" lang=\"ru\">на<span class=\"rohk\">у</span>чная пробл<span class=\"rohk\">е</span>ма</span><br> <span id=\"x43638_22_n\" class=\"n x_n n\" lang=\"et\">elu tekkimise probleem</span> <span id=\"x43638_23_qn\" class=\"qn x_qn qn\" lang=\"ru\">пробл<span class=\"rohk\">е</span>ма возникнов<span class=\"rohk\">е</span>ния ж<span class=\"rohk\">и</span>зни</span><br> <span id=\"x43638_24_n\" class=\"n x_n n\" lang=\"et\">arstiteaduse aktuaalseim probleem</span> <span id=\"x43638_25_qn\" class=\"qn x_qn qn\" lang=\"ru\">с<span class=\"rohk\">а</span>мая акту<span class=\"rohk\">а</span>льная пробл<span class=\"rohk\">е</span>ма в <span class=\"rohk\">о</span>бласти медиц<span class=\"rohk\">и</span>ны ~ медиц<span class=\"rohk\">и</span>нской на<span class=\"rohk\">у</span>ки</span><br> <span id=\"x43638_26_n\" class=\"n x_n n\" lang=\"et\">arutati tööhõive probleeme</span> <span id=\"x43638_27_qn\" class=\"qn x_qn qn\" lang=\"ru\">обсужд<span class=\"rohk\">а</span>лись пробл<span class=\"rohk\">е</span>мы ~ вопр<span class=\"rohk\">о</span>сы трудов<span class=\"rohk\">о</span>й з<span class=\"rohk\">а</span>нятости</span><br> <span id=\"x43638_28_n\" class=\"n x_n n\" lang=\"et\">millise probleemi kallal see labor töötab?</span> <span id=\"x43638_29_qn\" class=\"qn x_qn qn\" lang=\"ru\">над как<span class=\"rohk\">о</span>й пробл<span class=\"rohk\">е</span>мой <span class=\"rohk\">э</span>та лаборат<span class=\"rohk\">о</span>рия раб<span class=\"rohk\">о</span>тает?</span><br> <span id=\"x43638_30_n\" class=\"n x_n n\" lang=\"et\">selle eetilise probleemi pead lahendama ise</span> <span id=\"x43638_31_qn\" class=\"qn x_qn qn\" lang=\"ru\"><span class=\"rohk\">э</span>ту эт<span class=\"rohk\">и</span>ческую пробл<span class=\"rohk\">е</span>му ты д<span class=\"rohk\">о</span>лжен [раз]реш<span class=\"rohk\">и</span>ть сам</span><br> <span id=\"x43638_32_n\" class=\"n x_n n\" lang=\"et\">õppimine pole talle mingi probleem</span> <span id=\"x43638_33_qn\" class=\"qn x_qn qn\" lang=\"ru\">у нег<span class=\"rohk\">о</span> нет пробл<span class=\"rohk\">е</span>м с уч<span class=\"rohk\">ё</span>бой</span><br> <span id=\"x43638_34_n\" class=\"n x_n n\" lang=\"et\">meil pole rahaga probleeme</span> <span id=\"x43638_35_qn\" class=\"qn x_qn qn\" lang=\"ru\">д<span class=\"rohk\">е</span>ньги для н<span class=\"rohk\">а</span>с не пробл<span class=\"rohk\">е</span>ма</span> /  <span id=\"x43638_36_qn\" class=\"qn x_qn qn\" lang=\"ru\">у н<span class=\"rohk\">а</span>с нет пробл<span class=\"rohk\">е</span>м с деньг<span class=\"rohk\">а</span>ми</span><br> <span id=\"x43638_37_n\" class=\"n x_n n\" lang=\"et\">ta teeb igast asjast probleemi</span> <span id=\"x43638_38_qn\" class=\"qn x_qn qn\" lang=\"ru\">он д<span class=\"rohk\">е</span>лает из всег<span class=\"rohk\">о</span> пробл<span class=\"rohk\">е</span>му</span></p></div>	"
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(rawHTML))
	got, isMainArticle := GetSingleArticleWithDirection("probleem", doc.Nodes[0], "est-rus")
	expected := "<b>probleem</b><i> (probl'eem probleemi probl'eemi probl'eemi, probl'eemi[de probl'eemi[sid ~ probl'eem/e) </i>\r\n" +
		"проблема\r\n" +
		"вопрос"

	if got != expected {
		t.Errorf("bad card:\r\n %s\r\n%s", got, expected)
	}

	if isMainArticle != true {
		t.Errorf("article with grammar form should be the main")
	}
}

func TestFiltering(t *testing.T) {
	isMatch := IsMatchingArticle("antud", "'and[ma 'and[a anna[b 'an[tud, 'and[is 'and[ke")
	isGarbage := IsMatchingArticle("xyz", "'and[ma 'and[a anna[b 'an[tud, 'and[is 'and[ke")
	if isMatch != true && isGarbage != false {
		t.Error("IsMatchingArticle returns wrong result")
	}
}

func TestGetArticles(t *testing.T) {
	if os.Getenv("CI") != "" {
		t.Skip("Skipping testing in CI environment")
	}
	articles := GetArticles("laps")
	if len(articles) != 2 {
		t.Errorf("invalid number of articles: %d, expected 2", len(articles))
	}
}
