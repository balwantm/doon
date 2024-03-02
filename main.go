package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"regexp"
	"strings"
)

type Item struct {
	Name        string   `json:"name"`
	Price       string   `json:"price"`
	Description string   `json:"description"`
	ImageLinks  []string `json:"imageLinks"`
	Category    string   `json:"category"`
	Link        string   `json:"link"`
	Index       int
}

func main() {
	// Read data from data.json
	data, err := os.ReadFile("./data.json")
	if err != nil {
		fmt.Println("Error reading data.json:", err)
		return
	}

	// Unmarshal JSON data
	var items []Item
	err = json.Unmarshal(data, &items)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return
	}

	// Generate HTML for each item
	itemsHTML := generateItemsHTML(items)

	// Read index.html
	indexHTML, err := os.ReadFile("dist/index.html")
	if err != nil {
		fmt.Println("Error reading index.html:", err)
		return
	}

	// Replace <!-- INJECT_TABLE_HERE --> with the generated items HTML
	newIndexHTML := strings.Replace(string(indexHTML), "<!-- INJECT_TABLE_HERE -->", itemsHTML, 1)

	// Write the updated content back to index.html
	err = os.WriteFile("dist/index.html", []byte(newIndexHTML), 0644)
	if err != nil {
		fmt.Println("Error writing to index.html:", err)
		return
	}

	fmt.Println("Items successfully injected into index.html")

	// Create name+price.html for each item
	for i, item := range items {
		item.Index = i + 1 // Assuming Index is a field in your item struct

		err := createNamePriceFile(fmt.Sprintf("dist/%d.html", item.Index), newIndexHTML, item)
		if err != nil {
			fmt.Printf("Error creating %d.html: %v\n", item.Index, err)
		} else {
			fmt.Printf("%d.html successfully created\n", item.Index)
		}
	}
}

func generateItemsHTML(items []Item) string {
	// Define HTML template for the items
	itemTemplate := `
	<div class="col-sm-6 col-md-4 col-lg-3 p-b-35 isotope-item {{.Category}}">
		<div class="block2">
	<a href="{{.Index}}.html">
			<div class="block2-pic hov-img0">
				<img src="{{index .ImageLinks 0}}">
</a>
				
			</div>

			<div class="block2-txt flex-w flex-t p-t-14">
				<div class="block2-txt-child1 flex-col-l">
					<a href="{{.Index}}.html"
	class="stext-104 cl4 hov-cl1 trans-04 js-name-b2 p-b-6">
						{{.Name}}
					</a>

					<span class="stext-105 cl3">
						{{.Price}}
					</span>
				</div>
			</div>
		</div>
	</div>
	`

	// Create a template from the HTML
	tmpl, err := template.New("items").Parse(itemTemplate)
	if err != nil {
		fmt.Println("Error parsing HTML template:", err)
		return ""
	}

	// Create a buffer to store the rendered HTML
	var result strings.Builder

	// Execute the template for each item and write to the buffer
	for i, item := range items {
		item.Index = i + 1
		err := tmpl.Execute(&result, item)
		if err != nil {
			fmt.Println("Error executing template:", err)
			return ""
		}
	}

	return result.String()
}

func createNamePriceFile(filename string, newIndexHTML string, item Item) error {
	// Find the indices for header and footer
	headerEndIndex := strings.Index(newIndexHTML, "</header>") + len("</header>")
	footerStartIndex := strings.Index(newIndexHTML, "<footer>")

	// Extract content from start to </header> and <footer> to end
	headerContent := newIndexHTML[:headerEndIndex]
	footerContent := newIndexHTML[footerStartIndex:]

	// Create a template for the item's name
	nameTemplate := `<section class="sec-product-detail bg0 p-t-65 p-b-60">
		<div class="container">
			<div class="row">
				<div class="col-md-6 col-lg-7 p-b-30">
					<div class="p-l-25 p-r-30 p-lr-0-lg">
						<div class="wrap-slick3 flex-sb flex-w">
							<div class="wrap-slick3-dots"></div>
							<div class="wrap-slick3-arrows flex-sb-m flex-w"></div>

							<div class="slick3 gallery-lb">
	{{range .ImageLinks}}
								<div class="item-slick3" data-thumb="{{.}}">
									<div class="wrap-pic-w pos-relative">
										<img src="{{.}}" alt="IMG-PRODUCT">

										<a class="flex-c-m size-108 how-pos1 bor0 fs-16 cl10 bg0 hov-btn3 trans-04" href="{{.}}">
											<i class="fa fa-expand"></i>
										</a>
									</div>
								</div>
	{{end}}

								

								
							</div>
						</div>
					</div>
				</div>
					
				<div class="col-md-6 col-lg-5 p-b-30">
					<div class="p-r-50 p-t-5 p-lr-0-lg">
						<h4 class="mtext-105 cl2 js-name-detail p-b-14">
	{{.Name}}
						</h4>

						<span class="mtext-106 cl2">
	{{.Price}}
						</span>

						<p class="stext-102 cl3 p-t-23">
	{{.Description}}
						</p>
						
						<!--  -->
						<div class="p-t-33">
							<div class="flex-w flex-r-m p-b-10">
								

								<div class="size-204 respon6-next">
									
								</div>
							</div>

							<div class="flex-w flex-r-m p-b-10">
								

								<div class="size-204 respon6-next">
									
								</div>
							</div>

							<div class="flex-w flex-r-m p-b-10">
								<div class="size-204 flex-w flex-m respon6-next">
									
									<a href={{.Link}}>
									<button class="flex-c-m stext-101 cl0 size-101 bg1 bor1 hov-btn1 p-lr-15 trans-04 js-addcart-detail">
									Buy Now	
									</button></a>
								</div>
							</div>	
						</div>

						
					</div>
				</div>
			</div>

			
		</div>

		
	</section>

`
	tmpl, err := template.New("name").Parse(nameTemplate)
	if err != nil {
		return err
	}

	// Execute the template for the item
	var nameContent strings.Builder
	err = tmpl.Execute(&nameContent, item)
	if err != nil {
		return err
	}

	// Combine header, name content, and footer content
	finalContent := headerContent + nameContent.String() + footerContent

	// Write to dist/name+price.html
	err = os.WriteFile(filename, []byte(finalContent), 0644)
	if err != nil {
		return err
	}

	return nil
}

// sanitizeFileName replaces spaces with underscores and removes other invalid characters from a filename.
func sanitizeFileName(name string) string {
	invalidChars := regexp.MustCompile(`[^\w\d]+`)
	return invalidChars.ReplaceAllString(name, "_")
}
