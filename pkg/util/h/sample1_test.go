package h_test

import (
	"bytes"
	"net/http"
	"testing"

	"jdy/pkg/i18n"
	"jdy/pkg/lu"
	"jdy/pkg/util/h"
	"jdy/pkg/util/must"
)

var StaticURLPrefix = "/tsinghua"

func urlPrefix(intl *i18n.Context) string {
	switch intl.Language() {
	case i18n.ZhCN:
		return StaticURLPrefix
	case i18n.EnUS:
		return StaticURLPrefix + "/en"
	}
	panic("unsupported language" + intl.LangCode())
}

func homeURLPath(intl *i18n.Context) string            { return urlPrefix(intl) + "/" }
func researchURLPath(intl *i18n.Context) string        { return urlPrefix(intl) + "/research/" }
func industryURLPath(intl *i18n.Context) string        { return urlPrefix(intl) + "/industry/" }
func serviceURLPath(intl *i18n.Context) string         { return urlPrefix(intl) + "/service/" }
func campusLifeURLPath(intl *i18n.Context) string      { return urlPrefix(intl) + "/campus_life/" }
func tutorialsURLPath(intl *i18n.Context) string       { return urlPrefix(intl) + "/tutorials/" }
func tutorialContentURLPath(intl *i18n.Context) string { return urlPrefix(intl) + "/tutorial_content/" }
func courseURLPath(intl *i18n.Context) string          { return urlPrefix(intl) + "/course/" }
func programURLPath(intl *i18n.Context) string         { return urlPrefix(intl) + "/program/" }
func programContentURLPath(intl *i18n.Context) string  { return urlPrefix(intl) + "/program_content/" }
func opportunitiesURLPath(intl *i18n.Context) string   { return urlPrefix(intl) + "/opportunities/" }
func newsURLPath(intl *i18n.Context) string            { return urlPrefix(intl) + "/news/" }
func newsContentURLPath(intl *i18n.Context) string     { return urlPrefix(intl) + "/news_content/" }

type programURL struct {
	q            string
	programType  string
	degreeType   string
	locationType string
	targetType   string
	page         string
}

func newProgramURLFromParam(p *lu.Param) (programURL, error) {
	u := programURL{
		programType:  "all",
		degreeType:   "all",
		locationType: "all",
		targetType:   "all",
		page:         "1",
	}
	return u, p.
		Optional("programType", &u.programType).
		Optional("degreeType", &u.degreeType).
		Optional("locationType", &u.locationType).
		Optional("targetType", &u.targetType).
		Optional("page", &u.page).
		Optional("q", &u.q).Error()
}

func (u programURL) String(intl *i18n.Context) string {
	return h.URL(urlPrefix(intl)+"/program/", h.Params{
		{"programType", u.programType},
		{"degreeType", u.degreeType},
		{"locationType", u.locationType},
		{"targetType", u.targetType},
		{"page", u.page},
	})
}

func createProgramParse(intl *i18n.Context, locationTypeP string, programTypeP string, degreeTypeP string, targetTypeP string) string {
	newURL := programURL{
		programType:  "all",
		degreeType:   "all",
		locationType: "all",
		targetType:   "all",
		page:         "1",
	}
	return mCreatProgramParse(intl, newURL, locationTypeP, programTypeP, degreeTypeP, targetTypeP, "1")
}

func mCreatProgramParse(intl *i18n.Context, url programURL, locationTypeP string, programTypeP string, degreeTypeP string, targetTypeP string, page string) string {
	newURL := url
	if programTypeP != "" {
		newURL.programType = programTypeP
	}
	if degreeTypeP != "" {
		newURL.degreeType = degreeTypeP
	}
	if locationTypeP != "" {
		newURL.locationType = locationTypeP
	}
	if targetTypeP != "" {
		newURL.targetType = targetTypeP
	}
	if page != "" {
		newURL.page = page
	}
	return newURL.String(intl)
}

type tutorialsURL struct {
	tutorialType string
}

func newTutorialsURLFromParam(p *lu.Param) (tutorialsURL, error) {
	u := tutorialsURL{
		tutorialType: "all",
	}
	return u, p.
		Optional("tutorialType", &u.tutorialType).Error()
}

func (u tutorialsURL) String(intl *i18n.Context) string {
	return h.URL(urlPrefix(intl)+"/tutorials/", h.Params{
		{"tutorialType", u.tutorialType},
	})
}

func creatTutorialsParse(intl *i18n.Context, tutorialType string) string {
	newURL := tutorialsURL{
		tutorialType: "all",
	}
	if tutorialType != "" {
		newURL.tutorialType = tutorialType
	}
	return newURL.String(intl)
}

func HTMLBase(title string, head, body h.N) h.N {
	return h.H("html",
		h.H("head",
			h.H("title", title),
			h.H("meta", h.Attr("charset", "utf-8")),
			h.MetaNameContent("HandheldFriendly", "True"),
			h.MetaNameContent("MobileOptimized", "320"),
			h.MetaNameContent("mobile-web-app-capable", "yes"),
			h.MetaNameContent("viewport", "width=device-width, initial-scale=1, user-scalable=no"),
			h.H("meta", h.Attrs{{"http-equiv", "cleartype"}, {"content", "on"}}),
			head,
		),
		h.H("body", body),
	)
}

func staticURL(s string) string { return s }

func base(intl *i18n.Context, title string, content h.N) h.N {
	head := h.L(
		h.H("link", h.Attrs{{"href", staticURL("logo-icon.png")}, {"rel", "shortcut icon"}}),
		h.H("link", h.Attrs{{"href", "//dn-applysquare-lib.qbox.me/bootstrap/3.3.6/css/bootstrap.min.css"}, {"rel", "stylesheet"}}),
		h.H("link", h.Attrs{{"href", "/dist/tsinghua.css"}, {"rel", "stylesheet"}}),
	)

	body := h.L(
		content,
		h.H(".footer",
			h.H(".related-links",
				h.H(".container",
					h.H(".row",
						h.H(".col-md-2",
							h.H("ul.links-list",
								h.H("li.links-title",
									h.H("span"),
									h.H("a", intl.S("学校概况")),
								),
								h.H("li",
									h.H("a", intl.S("校长致辞")),
								),
								h.H("li",
									h.H("a", intl.S("学校沿革")),
								),
								h.H("li",
									h.H("a", intl.S("历届领导")),
								),
								h.H("li",
									h.H("a", intl.S("现任领导")),
								),
								h.H("li",
									h.H("a", intl.S("组织机构")),
								),
								h.H("li",
									h.H("a", intl.S("统计资料")),
								),
							),
						),
						h.H(".col-md-2",
							h.H("ul.links-list",
								h.H("li.links-title",
									h.H("span"),
									h.H("a", intl.S("院系设置")),
								),
							),
							h.H("ul.links-list",
								h.H("li.links-title",
									h.H("span"),
									h.H("a", intl.S("师资队伍")),
								),
								h.H("li",
									h.H("a", intl.S("概况")),
								),
								h.H("li",
									h.H("a", intl.S("杰出人才")),
								),
							),
						),
						h.H(".col-md-2",
							h.H("ul.links-list",
								h.H("li.links-title",
									h.H("span"),
									h.H("a", intl.S("教育教学")),
								),
								h.H("li",
									h.H("a", intl.S("本科生教育")),
								),
								h.H("li",
									h.H("a", intl.S("研究生教育")),
								),
								h.H("li",
									h.H("a", intl.S("留学生教育")),
								),
								h.H("li",
									h.H("a", intl.S("继续教育")),
								),
							),
						),
						h.H(".col-md-2",
							h.H("ul.links-list",
								h.H("li.links-title",
									h.H("span"),
									h.H("a", intl.S("科学研究")),
								),
								h.H("li",
									h.H("a", intl.S("科研项目")),
								),
								h.H("li",
									h.H("a", intl.S("科研机构")),
								),
								h.H("li",
									h.H("a", intl.S("科研合作")),
								),
								h.H("li",
									h.H("a", intl.S("科研成果与知识产权")),
								),
								h.H("li",
									h.H("a", intl.S("学术交流")),
								),
							),
						),
						h.H(".col-md-2",
							h.H("ul.links-list",
								h.H("li.links-title",
									h.H("span"),
									h.H("a", intl.S("招生就业")),
								),
								h.H("li",
									h.H("a", intl.S("本科生招生")),
								),
								h.H("li",
									h.H("a", intl.S("研究生招生")),
								),
								h.H("li",
									h.H("a", intl.S("留学生招生")),
								),
								h.H("li",
									h.H("a", intl.S("学生职业发展")),
								),
							),
						),
						h.H(".col-md-2",
							h.H("ul.links-list",
								h.H("li.links-title",
									h.H("span"),
									h.H("a", intl.S("走进清华")),
								),
								h.H("li",
									h.H("a", intl.S("校园生活")),
								),
								h.H("li",
									h.H("a", intl.S("校园风光")),
								),
								h.H("li",
									h.H("a", intl.S("实用信息")),
								),
							),
						),
					),
				),
			),
			h.H(".record-information",
				h.H(".container",
					h.H(".row",
						h.H(".col-sm-6.col-md-4.col-lg-3.sp-bottom-2x", intl.S("电话查号台：010-62793001")),
						h.H(".col-sm-6.col-md-4.col-lg-3.sp-bottom-2x", intl.S("管理员信箱：xinxiban@tsinghua.edu.cn")),
						h.H(".col-sm-6.col-md-4.col-lg-3.sp-bottom-2x", intl.S("地址：北京市海淀区清华大学")),
						h.H(".col-sm-6.col-md-4.col-lg-3.sp-bottom-2x", intl.S("京公网安备 110402430053 号")),
						h.H(".col-sm-6.col-md-4.col-lg-3.sp-bottom-2x", intl.S("版权所有 © 清华大学")),
					),
				),
			),
		),

		h.H("script", h.Attr("src", "//dn-applysquare-lib.qbox.me/jquery/1.11.3/jquery.min.js")),
		h.H("script", h.Attr("src", "//dn-applysquare-lib.qbox.me/bootstrap/3.3.6/js/bootstrap.min.js")),
	)

	return HTMLBase(title, head, body)
}

func isAction(class string, nav string) string {
	if class == nav {
		return ".active"
	}
	return ""
}

func navInc(intl *i18n.Context, webNav, webSubNav string, req *http.Request) h.N {
	return h.H(".header."+webNav,
		h.H(".container",
			h.H(".logo.float-left",
				h.H("img", h.Attr("src", staticURL("logo-flag.png"))),
			),
			h.H(".info.float-left",
				h.H("p.name.ellipsis", intl.S("Tsinghua University International Education")),
				h.H("p.ellipsis", h.S("Tsinghua International Education")),
			),
			h.H(".link.slide.float-right.sp-top-2x.hidden-sm.hidden-xs",
				h.H("a", h.Href("#"), intl.S("ABOUT US")),
				h.H("span.color-peru", h.S("●")),
				h.H("a", h.Href("#"), intl.S("CONTACT US")),
				h.H("span.color-peru", h.S("●")),
				h.H("a"+isAction("zh-cn", intl.LangCode()), h.Href("#"), h.S("中文")),
				h.H("span", h.S("/")),
				h.H("a"+isAction("en-us", intl.LangCode()), h.Href("#"), h.S("English")),
			),
			h.H(".clear"),
		),
		h.H("nav.navbar.navbar-absolute-top.bs-docs-nav",
			h.H(".container",
				h.H(".navbar-header",
					h.H("button.navbar-toggle.collapsed", h.Attrs{{"type", "button"}, {"data-toggle", "collapse"}, {"data-target", "#navbar"}, {"aria-expanded", "false"}, {"aria-controls", "navbar"}},
						h.H("span.sr-only", intl.S("Toggle navigation")),
						h.H("span.icon-bar", h.S("")),
						h.H("span.icon-bar", h.S("")),
						h.H("span.icon-bar", h.S("")),
					),
				),
				h.H("#navbar.collapse.navbar-collapse",
					h.H("ul.nav.navbar-nav",
						h.H("li"+isAction("home", webNav),
							h.H("a", h.Href(homeURLPath(intl)), intl.S("Home")),
						),
						h.H("li.dropdown"+isAction("overview", webNav),
							h.H("a#drop1.dropdown-toggle",
								h.Attrs{{"data-toggle", "dropdown"}, {"role", "button"}, {"aria-haspopup", "true"}, {"aria-expanded", "false"}},
								intl.S("Overview"),
								h.H("span.caret"),
							),
							h.H("ul.dropdown-menu", h.Attr("aria-labelledby", "drop1"),
								h.H("li"+isAction("industry", webSubNav),
									h.H("a", h.Href(industryURLPath(intl)), intl.S("Global Vision")),
								),
								h.H("li"+isAction("industry", webSubNav),
									h.H("a", h.Href(industryURLPath(intl)), intl.S("International Degree Program")),
								),
								h.H("li"+isAction("research", webSubNav),
									h.H("a", h.Href(researchURLPath(intl)), intl.S("Study and Research Abroad")),
								),
								h.H("li"+isAction("service", webSubNav),
									h.H("a", h.Href(serviceURLPath(intl)), intl.S("Service and Internship")),
								),
								h.H("li"+isAction("campus_life", webSubNav),
									h.H("a", h.Href(campusLifeURLPath(intl)), intl.S("Community and Campus Life")),
								),
							),
						),
						h.H("li.dropdown"+isAction("tutorials", webNav),
							h.H("a#drop2.dropdown-toggle",
								h.Attrs{{"data-toggle", "dropdown"}, {"role", "button"}, {"aria-haspopup", "true"}, {"aria-expanded", "false"}},
								intl.S("Tutorial"),
								h.H("span.caret"),
							),
							h.H("ul.dropdown-menu", h.Attr("aria-labelledby", "drop2"),
								h.H("li"+isAction("preparation", webSubNav),
									h.H("a", h.Href(creatTutorialsParse(intl, "preparation")), intl.S("Why Go Global")),
								),
								h.H("li"+isAction("improvement", webSubNav),
									h.H("a", h.Href(creatTutorialsParse(intl, "improvement")), intl.S("Are You Ready")),
								),
								h.H("li"+isAction("exploration", webSubNav),
									h.H("a", h.Href(creatTutorialsParse(intl, "exploration")), intl.S("Where to Find")),
								),
								h.H("li"+isAction("application", webSubNav),
									h.H("a", h.Href(creatTutorialsParse(intl, "application")), intl.S("How to Apply")),
								),
							),
						),
						h.H("li.dropdown"+isAction("study", webNav),
							h.H("a#drop3.dropdown-toggle",
								h.Attrs{{"data-toggle", "dropdown"}, {"role", "button"}, {"aria-haspopup", "true"}, {"aria-expanded", "false"}},
								intl.S("Learn"),
								h.H("span.caret"),
							),
							h.H("ul.dropdown-menu", h.Attr("aria-labelledby", "drop3"),
								h.H("li"+isAction("degree_program", webSubNav),
									h.H("a", h.Href(createProgramParse(intl, "", "degree_program", "", "")), intl.S("Degree Program")),
								),
								h.H("li"+isAction("exchange_program", webSubNav),
									h.H("a", h.Href(createProgramParse(intl, "", "exchange_program", "", "")), intl.S("Exchange Program")),
								),
								h.H("li"+isAction("summer_school", webSubNav),
									h.H("a", h.Href(createProgramParse(intl, "", "summer_school", "", "")), intl.S("Summer Program")),
								),
							),
						),
						h.H("li.dropdown"+isAction("research", webNav),
							h.H("a#drop3.dropdown-toggle",
								h.Attrs{{"data-toggle", "dropdown"}, {"role", "button"}, {"aria-haspopup", "true"}, {"aria-expanded", "false"}},
								intl.S("Research"),
								h.H("span.caret"),
							),
							h.H("ul.dropdown-menu", h.Attr("aria-labelledby", "drop3"),
								h.H("li"+isAction("government_sponsored", webSubNav),
									h.H("a", h.Href(createProgramParse(intl, "", "government_sponsored", "", "")), intl.S("Government-sponsored Research")),
								),
								h.H("li"+isAction("overseas_thesis", webSubNav),
									h.H("a", h.Href(createProgramParse(intl, "", "overseas_thesis", "", "")), intl.S("Thesis Project")),
								),
								h.H("li"+isAction("research_assistant", webSubNav),
									h.H("a", h.Href(createProgramParse(intl, "", "research_assistant", "", "")), intl.S("Research Internship")),
								),
								h.H("li"+isAction("visiting_scholar", webSubNav),
									h.H("a", h.Href(createProgramParse(intl, "", "visiting_scholar", "", "")), intl.S("Visiting Scholar")),
								),
							),
						),
						h.H("li.dropdown"+isAction("events", webNav),
							h.H("a#drop3.dropdown-toggle",
								h.Attrs{{"data-toggle", "dropdown"}, {"role", "button"}, {"aria-haspopup", "true"}, {"aria-expanded", "false"}},
								intl.S("Activities"),
								h.H("span.caret"),
							),
							h.H("ul.dropdown-menu", h.Attr("aria-labelledby", "drop3"),
								h.H("li"+isAction("conference", webSubNav),
									h.H("a", h.Href(createProgramParse(intl, "", "conference", "", "")), intl.S("Conference")),
								),
								h.H("li"+isAction("seminar,workshop,lecture_forum", webSubNav),
									h.H("a", h.Href(createProgramParse(intl, "", "seminar,workshop,lecture_forum", "", "")), intl.S("Seminar, Workshop and Lecture")),
								),
								h.H("li"+isAction("competition", webSubNav),
									h.H("a", h.Href(createProgramParse(intl, "", "competition", "", "")), intl.S("Competition")),
								),
								h.H("li"+isAction("internship", webSubNav),
									h.H("a", h.Href(createProgramParse(intl, "", "internship", "", "")), intl.S("Industry Internship")),
								),
								h.H("li"+isAction("cultural_exchange", webSubNav),
									h.H("a", h.Href(createProgramParse(intl, "", "cultural_exchange", "", "")), intl.S("Cultural Exchange")),
								),
								h.H("li"+isAction("international_volunteers", webSubNav),
									h.H("a", h.Href(createProgramParse(intl, "", "international_volunteers", "", "")), intl.S("International Volunteers")),
								),
							),
						),
					),
					h.H("form.navbar-form.navbar-right.hide", h.Attr("role", "search"),
						h.H(".form-group",
							h.H("input.form-control", h.Attrs{{"type", "text"}, {"placeholder", "Search"}}),
						),
					),
				),
			),
		),
	)
}

func bannerCarousel(intl *i18n.Context) h.N {
	return h.L(
		h.H(".row",
			h.H(".col-md-12.sp-bottom-2x",
				h.H("h2.slogan.center", intl.S("Rooting in China, Reaching for the World")),
			),
			h.H(".col-md-4",
				h.H(".menu",
					h.H("p.menu_title", intl.S("Global Competence")),
					h.H("p",
						h.H("a", h.Href(industryURLPath(intl)), intl.S("International Degree Program")),
					),
					h.H("p",
						h.H("a", h.Href(researchURLPath(intl)), intl.S("Study and Research Abroad")),
					),
					h.H("p",
						h.H("a", h.Href(serviceURLPath(intl)), intl.S("Service and Internship")),
					),
					h.H("p",
						h.H("a", h.Href(campusLifeURLPath(intl)), intl.S("Community and Campus Life")),
					),
				),
			),
			h.H(".col-md-4",
				h.H(".menu",
					h.H("p.menu_title", intl.S("Preparation")),
					h.H("p",
						h.H("a", h.Href(creatTutorialsParse(intl, "preparation")), intl.S("Why Go Global")),
					),
					h.H("p",
						h.H("a", h.Href(creatTutorialsParse(intl, "improvement")), intl.S("Are You Ready")),
					),
					h.H("p",
						h.H("a", h.Href(creatTutorialsParse(intl, "exploration")), intl.S("Where to Find")),
					),
					h.H("p",
						h.H("a", h.Href(creatTutorialsParse(intl, "application")), intl.S("How to Apply")),
					),
				),
			),
			h.H(".col-md-4",
				h.H(".menu",
					h.H("p.menu_title", intl.S("Explore")),
					h.H("p",
						h.H("a", h.Href(createProgramParse(intl, "", "", "undergraduate_students", "")), intl.S("For Undergraduates")),
					),
					h.H("p",
						h.H("a", h.Href(createProgramParse(intl, "", "", "master", "")), intl.S("For Postgraduates")),
					),
					h.H("p",
						h.H("a", h.Href(createProgramParse(intl, "", "", "phd", "")), intl.S("For Doctoral Students")),
					),
					h.H("p",
						h.H("a", h.Href(createProgramParse(intl, "", "", "", "current_foreign")), intl.S("Current International Students")),
					),
				),
			),
		),
	)
}

func homeNode(intl *i18n.Context) h.N {
	return base(intl, intl.S("Home Page"), h.L(
		navInc(intl, "home", "", nil),
		h.H("#carousel-example-generic.carousel.slide.hidden-sm.hidden-xs", h.Attr("data-ride", "carousel"),
			h.H("ol.carousel-indicators",
				h.H("li.active", h.Attrs{{"data-target", "#carousel-example-generic"}, {"data-slide-to", "0"}}),
				h.H("li", h.Attrs{{"data-target", "#carousel-example-generic"}, {"data-slide-to", "1"}}),
				h.H("li", h.Attrs{{"data-target", "#carousel-example-generic"}, {"data-slide-to", "2"}}),
			),
			h.H(".carousel-inner", h.Attr("role", "listbox"),
				h.H(".item.active",
					h.H("img", h.Attr("style", "background-image: url("+staticURL("banner-0.png")+")")),
				),
				h.H(".item",
					h.H("img", h.Attr("style", "background-image: url("+staticURL("banner-1.png")+")")),
				),
				h.H(".item",
					h.H("img", h.Attr("style", "background-image: url("+staticURL("banner-2.png")+")")),
				),
			),
			h.H(".carousel-caption",
				bannerCarousel(intl),
			),
			h.H("a.left.carousel-control", h.Attrs{{"href", "#carousel-example-generic"}, {"role", "button"}, {"data-slide", "prev"}},
				h.H("span.glyphicon.glyphicon-chevron-left", h.Attr("aria-hidden", "true")),
				h.H("span.sr-only", h.S("Previous")),
			),
			h.H("a.right.carousel-control", h.Attrs{{"href", "#carousel-example-generic"}, {"role", "button"}, {"data-slide", "next"}},
				h.H("span.glyphicon.glyphicon-chevron-right", h.Attr("aria-hidden", "true")),
				h.H("span.sr-only", h.S("Next")),
			),
		),
		h.H(".container.academic-programs-num",
			h.H(".row",
				h.H(".col-md-3",
					h.H("p.num", h.S("15,000")),
					h.H("p", intl.S("覆盖学生人数")),
				),
				h.H(".division"),
				h.H(".col-md-3",
					h.H("p.num", h.S("3000+")),
					h.H("p", intl.S("国际研究项目")),
				),
				h.H(".division"),
				h.H(".col-md-3",
					h.H("p.num", h.S("368")),
					h.H("p", intl.S("国际合作院校")),
				),
				h.H(".division"),
				h.H(".col-md-3",
					h.H("p.num", h.S("126")),
					h.H("p", intl.S("合作国家")),
				),
			),
		),
		h.H(".title.sp-bottom-3x",
			h.H("span.tab", intl.S("新闻动态")),
			h.H("span.icon"),
		),

		h.H(".container.trends.sp-top-3x",
			h.H(".row.sp-top-3x",
				h.H(".col-md-6",
					h.H("p.center",
						h.H("span.color-peru.size-16", intl.S("新闻")),
						h.H("span.color-blue", h.S("|")),
						h.H("a.color-blue", h.Href(newsURLPath(intl)), intl.S("更多新闻")),
					),
					h.H("hr.saddlebrown.sp-top.sp-bottom-2x"),
					h.H("img.sp-bottom-2x", h.Attr("src", staticURL("news-1.png"))),
					h.H("a.sp-top-2x",
						h.Href(newsContentURLPath(intl)),
						h.H("p.news_title.size-18", intl.S("美国前财长劳伦斯·萨默斯受聘清华大学杰出访问教授")),
						h.H("p", intl.S("11月4日，清华大学杰出访问教授聘任仪式在清华大学音乐厅举行。美国前财政部长、哈佛大学前校长劳伦斯·萨默斯（Lawrence H. Summers）受聘为清华大学杰出访问教授,未来将正式加入清华大学苏世民学者项目的教学工作。清华大学校长邱勇出席仪式...")),
					),
					h.H("a.sp-top-3x",
						h.Href(newsContentURLPath(intl)),
						h.H("p.news_title.size-18", intl.S("美国前财长劳伦斯·萨默斯受聘清华大学杰出访问教授")),
						h.H("p", intl.S("11月4日，清华大学杰出访问教授聘任仪式在清华大学音乐厅举行。美国前财政部长、哈佛大学前校长劳伦斯·萨默斯（Lawrence H. Summers）受聘为清华大学杰出访问教授,未来将正式加入清华大学苏世民学者项目的教学工作。清华大学校长邱勇出席仪式...")),
					),
					h.H("a.sp-top-3x",
						h.Href(newsContentURLPath(intl)),
						h.H("p.news_title.size-18", intl.S("美国前财长劳伦斯·萨默斯受聘清华大学杰出访问教授")),
						h.H("p", intl.S("11月4日，清华大学杰出访问教授聘任仪式在清华大学音乐厅举行。美国前财政部长、哈佛大学前校长劳伦斯·萨默斯（Lawrence H. Summers）受聘为清华大学杰出访问教授,未来将正式加入清华大学苏世民学者项目的教学工作。清华大学校长邱勇出席仪式...")),
					),
				),
				h.H(".col-md-6",
					h.H("p.center",
						h.H("span.color-peru.size-16", intl.S("通知公告")),
						h.H("span.color-blue", h.S("|")),
						h.H("a.color-blue", h.Href(newsURLPath(intl)), intl.S("更多通知")),
					),
					h.H("hr.saddlebrown.sp-top.sp-bottom-2x"),
					h.H(".notify-item",
						h.H(".date",
							h.H(".day", h.S("03-10")),
							h.H(".year", h.S("2016")),
						),
						h.H(".text",
							h.H("p",
								h.H("a", h.Href(newsContentURLPath(intl)), intl.S("关于报考清华大学2016年博士生资格审查的通知")),
							),
							h.H("p",
								h.H("a", h.Href(newsContentURLPath(intl)), intl.S("清华大学五道口金融学院2016年金融学直博生推免项目招生说明")),
							),
							h.H("p",
								h.H("a", h.Href(newsContentURLPath(intl)), intl.S("关于保留入学资格学生2016年返校的通知")),
							),
							h.H("p",
								h.H("a", h.Href(newsContentURLPath(intl)), intl.S("探索顶尖科学需聚世界之力－引力波探测全球超过千名科学家参与")),
							),
						),
					),
					h.H(".notify-item",
						h.H(".date",
							h.H(".day", h.S("03-10")),
							h.H(".year", h.S("2016")),
						),
						h.H(".text",
							h.H("p",
								h.H("a", h.Href(newsContentURLPath(intl)), intl.S("清华大学2016年博士研究生入学考试准考证打印及交通餐饮提示")),
							),
							h.H("p",
								h.H("a", h.Href(newsContentURLPath(intl)), intl.S("清华大学柳斌杰：加大力度培养出版产业职业经理")),
							),
							h.H("p",
								h.H("a", h.Href(newsContentURLPath(intl)), intl.S("清华大学苏世民学者项目2016年招生报名申请系统开放")),
							),
						),
					),
					h.H(".notify-item",
						h.H(".date",
							h.H(".day", h.S("03-10")),
							h.H(".year", h.S("2016")),
						),
						h.H(".text",
							h.H("p",
								h.H("a", h.Href(newsContentURLPath(intl)), intl.S("2015年清华大学苏世民学者项目夏令营报名通知")),
							),
							h.H("p",
								h.H("a", h.Href(newsContentURLPath(intl)), intl.S("关于报考清华大学2016年博士生资格审查的通知")),
							),
							h.H("p",
								h.H("a", h.Href(newsContentURLPath(intl)), intl.S("清华大学2016年博士研究生入学考试准考证打印及交通餐饮提示")),
							),
							h.H("p",
								h.H("a", h.Href(newsContentURLPath(intl)), intl.S("关于报考清华大学2016年博士生资格审查的通知")),
							),
						),
					),
					h.H(".notify-item",
						h.H(".date",
							h.H(".day", h.S("03-10")),
							h.H(".year", h.S("2016")),
						),
						h.H(".text",
							h.H("p",
								h.H("a", h.Href(newsContentURLPath(intl)), intl.S("2015年清华大学苏世民学者项目夏令营报名通知")),
							),
							h.H("p",
								h.H("a", h.Href(newsContentURLPath(intl)), intl.S("关于报考清华大学2016年博士生资格审查的通知")),
							),
							h.H("p",
								h.H("a", h.Href(newsContentURLPath(intl)), intl.S("清华大学2016年博士研究生入学考试准考证打印及交通餐饮提示")),
							),
							h.H("p",
								h.H("a", h.Href(newsContentURLPath(intl)), intl.S("关于报考清华大学2016年博士生资格审查的通知")),
							),
						),
					),
				),
			),
		),
		h.H(".title.sp-bottom-3x",
			h.H("span.tab", intl.S("项目展示")),
			h.H("span.icon"),
		),

		h.H(".container.trends",
			h.H(".row",
				h.H(".col-md-4",
					h.H("hr.saddlebrown"),
					h.H("img.sp-bottom-2x", h.Attr("src", staticURL("news-1.png"))),
					h.H("a.sp-top-2x",
						h.Href(newsContentURLPath(intl)),
						h.H("p.news_title.size-16", intl.S("美国前财长劳伦斯·萨默斯受聘清华大学杰出访问教授")),
						h.H("p.sp-top-2x", intl.S("11月4日，清华大学杰出访问教授聘任仪式在清华大学音乐厅举行。美国前财政部长、哈佛大学前校长劳伦斯·萨默斯（Lawrence H. Summers）受聘为清华大学杰出访问教授,未来将正式加入清华大学苏世民学者项目的教学工作...")),
					),
				),
				h.H(".col-md-4",
					h.H("hr.saddlebrown"),
					h.H("img.sp-bottom-2x", h.Attr("src", staticURL("news-2.png"))),
					h.H("a.sp-top-2x",
						h.Href(newsContentURLPath(intl)),
						h.H("p.news_title.size-16", intl.S("美国前财长劳伦斯·萨默斯受聘清华大学杰出访问教授")),
						h.H("p.sp-top-2x", intl.S("11月4日，清华大学杰出访问教授聘任仪式在清华大学音乐厅举行。美国前财政部长、哈佛大学前校长劳伦斯·萨默斯（Lawrence H. Summers）受聘为清华大学杰出访问教授,未来将正式加入清华大学苏世民学者项目的教学工作...")),
					),
				),
				h.H(".col-md-4",
					h.H("hr.saddlebrown"),
					h.H("img.sp-bottom-2x", h.Attr("src", staticURL("news-1.png"))),
					h.H("a.sp-top-2x",
						h.Href(newsContentURLPath(intl)),
						h.H("p.news_title.size-16", intl.S("美国前财长劳伦斯·萨默斯受聘清华大学杰出访问教授")),
						h.H("p.sp-top-2x", intl.S("11月4日，清华大学杰出访问教授聘任仪式在清华大学音乐厅举行。美国前财政部长、哈佛大学前校长劳伦斯·萨默斯（Lawrence H. Summers）受聘为清华大学杰出访问教授,未来将正式加入清华大学苏世民学者项目的教学工作...")),
					),
				),
			),
		),
		h.H(".title.sp-bottom-3x",
			h.H("span.tab", intl.S("图文专题")),
			h.H("span.icon"),
		),
		h.H(".container.sp-top-3x.sp-bottom-3x",
			h.H(".row.sp-top-3x",
				h.H(".col-md-4",
					h.H(".stories-item", h.Attr("style", "background-image: url("+staticURL("stories-1.png")+")"),
						h.H(".text.xviolet", intl.S("清华大学 - 美国北卡罗莱纳大学全球供应链领袖双硕士学位项目")),
					),
				),
				h.H(".col-md-4",
					h.H(".stories-item", h.Attr("style", "background-image: url("+staticURL("stories-2.png")+")"),
						h.H(".text.xviolet", intl.S("清华大学 - 美国北卡罗莱纳大学全球供应链领袖双硕士学位项目"))),
				),
				h.H(".col-md-4",
					h.H(".stories-item", h.Attr("style", "background-image: url("+staticURL("stories-3.png")+")"),
						h.H(".text.xviolet", intl.S("清华大学 - 美国北卡罗莱纳大学全球供应链领袖双硕士学位项目"))),
				),
			),
		),
	))
}

func TestHome(t *testing.T) {
	t.Skip()
	var w bytes.Buffer
	err := homeNode(i18n.EnContext).ToHTML(&w)
	if err != nil {
		t.Error(err)
	}
	t.Error(string(w.Bytes()))
}

var k = 0

func BenchmarkHome(b *testing.B) {
	for n := 0; n < b.N; n++ {
		var w bytes.Buffer
		err := homeNode(i18n.EnContext).ToHTML(&w)
		must.Must(err)
		k += w.Len()
	}
}
