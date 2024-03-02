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
	Name        string `json:"name"`
	Price       string `json:"price"`
	Description string `json:"description"`
	ImageLink   string `json:"imageLinks"`
	Category    string `json:"category"`
	Link        string `json:"link"`
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
	<div class="w-full md:w-1/3 xl:w-1/4 p-6 flex flex-col">
				<a href="{{.Index}}.html">
					<img class="hover:grow hover:shadow-lg"
						src="{{.ImageLink}}">
					<div class="pt-3 flex items-center justify-between">
						<p class="">{{.Name}}</p>
											</div>
					<p class="pt-1 text-gray-900">{{.Price}}</p>
				</a>
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
	nameTemplate := `
<main class="w-full flex flex-col lg:flex-row">
				<!-- Gallery -->
				<section
					class="h-fit flex-col gap-8 mt-4 sm:flex sm:flex-row sm:gap-4 sm:h-full sm:mt-6 sm:mx-2 md:gap-8 md:mx-4 lg:flex-col lg:mx-0 lg:mt-36"
				>
					<picture
						class="relative flex items-center bg-orange sm:bg-transparent"
					>
						
						<img
							src="{{.ImageLink}}"
							alt="sneaker"
							class="block sm:rounded-xl xl:w-[70%] xl:rounded-xl m-auto pointer-events-none transition duration-300 lg:w-3/4 lg:pointer-events-auto lg:cursor-pointer lg:hover:shadow-xl"
							id="hero"
						/>
						
					</picture>

									</section>

				<!-- Text -->
				<section
					class="w-full p-6 lg:mt-36 lg:pr-20 lg:py-10 2xl:pr-40 2xl:mt-40"
				>
										<h1
						class="text-very-dark m-4 font-bold text-3xl lg:text-4xl"
					>
	{{.Name}}	
					</h1>
					<p class="text-dark-grayish mb-6 text-base sm:text-lg">
	{{.Description}}	
					</p>

					<div
						class="flex items-center justify-between mb-6 sm:flex-col sm:items-start"
					>
						<div class="flex items-center gap-4">
							<h3
								class="text-very-dark font-bold text-3xl inline-block"
							>
	{{.Price}}	
							</h3>
													</div>
											</div>
<a href="{{.Link}}">
											<button
							class="w-full h-10 bg-orange py-2 flex items-center justify-center gap-4 text-xl rounded-lg font-bold text-light shadow-md shadow-orange hover:brightness-125 transition select-none"
							id="add-cart"
						>
							<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-cart2" viewBox="0 0 16 16">
  <path d="M0 2.5A.5.5 0 0 1 .5 2H2a.5.5 0 0 1 .485.379L2.89 4H14.5a.5.5 0 0 1 .485.621l-1.5 6A.5.5 0 0 1 13 11H4a.5.5 0 0 1-.485-.379L1.61 3H.5a.5.5 0 0 1-.5-.5M3.14 5l1.25 5h8.22l1.25-5zM5 13a1 1 0 1 0 0 2 1 1 0 0 0 0-2m-2 1a2 2 0 1 1 4 0 2 2 0 0 1-4 0m9-1a1 1 0 1 0 0 2 1 1 0 0 0 0-2m-2 1a2 2 0 1 1 4 0 2 2 0 0 1-4 0"/>
</svg>						Buy Now	
						</button></a>
					</div>
				</section>
			</main>

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
