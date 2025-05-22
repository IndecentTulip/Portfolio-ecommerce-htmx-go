package db_api
import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"
	"strings"

	wc "HtmxReactGolang/web_context"

	_ "github.com/mattn/go-sqlite3"
	_ "github.com/lib/pq"
)

func DeleteFromProducts(db *sql.DB, sessionID string ){
  deletequery := `DELETE FROM products WHERE id = $1`
  cartItems := SelectCart(db,sessionID)
  for _,item := range cartItems{
    fmt.Println("HELLO I AM TRYING TO DELETE: " + item.Product.Id)
    _, err := db.Exec(deletequery, item.Product.Id)
    if err != nil {
      log.Fatal(err)
    }
  }
}

func GetProductsList(db *sql.DB, offset int) ([]wc.Product, int) {
	query := `
    SELECT 
        COUNT(*) OVER() AS total,
        p.id, p.name, p.price, p.descript, p.quantity, p.image,
        STRING_AGG(pt.tagName, ',') AS tags
    FROM 
        products p
    LEFT JOIN 
        productTags pt ON p.id = pt.ProductId
    GROUP BY 
        p.id
    ORDER BY 
        p.price
    LIMIT 10
    OFFSET $1
`

	var ProductList []wc.Product

	rows, err := db.Query(query, offset)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var id, name, desc, tagsStr string
	var price, total, quantity int
	var imgByte []byte

	for rows.Next() {
		err := rows.Scan(&total, &id, &name, &price, &desc, &quantity, &imgByte, &tagsStr)
		if err != nil {
			log.Fatal(err)
		}

		imgStr := base64.StdEncoding.EncodeToString(imgByte)

		var tags []string
		if tagsStr != "" {
			tags = strings.Split(tagsStr, ",")
		}

		ProductList = append(ProductList, wc.Product{
			Id:       id,
			Name:     name,
			Price:    price,
			Desc:     desc,
			Quantity: quantity,
			Image:    imgStr,
			Tags:     tags,
		})
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return ProductList, total
}


func GetProduct(db *sql.DB, prodId string) wc.Product {
  query := `SELECT id, name, price, desc, quantity, image FROM products WHERE id = $1`

  var product wc.Product

  rows, err := db.Query(query, prodId)
  if err != nil {
    log.Fatal(err)
  }

  var id, name, desc string
  var price, quantity int
  var imgByte []byte

  for rows.Next() {
    err := rows.Scan(&id, &name, &price, &desc, &quantity, &imgByte)
    if err != nil {
        log.Fatal(err)
    }
    
    imgStr := base64.StdEncoding.EncodeToString(imgByte)
    product = wc.Product{Id: id, Name: name, Price: price, Desc: desc, Quantity: quantity, Image: imgStr}
  }

    defer rows.Close()

  return product

}

func GetProductListSearch(db *sql.DB, term string, offset int) ([]wc.Product, int) {
	var query string
	query = `
	SELECT 
			COUNT(*) OVER() AS total,
      p.id, p.name, p.price, p.descript, p.quantity, p.image,
			STRING_AGG(pt.tagName, ',') AS tags,
			ts_rank(search, websearch_to_tsquery('english', $1)) +
			ts_rank(search, websearch_to_tsquery('simple', $1)) AS rank
	FROM 
			products p
	LEFT JOIN 
			productTags pt ON p.id = pt.ProductId
	WHERE 
		p.search @@ websearch_to_tsquery('english', $1)
	OR
		p.search @@ websearch_to_tsquery('simple', $1)
	GROUP BY 
			p.id
	ORDER BY 
			rank desc, p.price
	LIMIT 10
	OFFSET $2
`
	var ProductList []wc.Product

	rows, err := db.Query(query, term, offset)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var id, name, desc, tagsStr string
	var price, total, quantity int
	var imgByte []byte
	var rank float64

	booll := true

	// TODO change a way we tell that first expression failed
	// TODO change query, it's so shit 
  if !booll {
		fmt.Println("$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$")
		query = `
		WITH search_terms AS (
				SELECT unnest(string_to_array($1, ' ')) AS term
		),
		tag_matches AS (
				SELECT
						pt.ProductId,
						-- Check if any tag matches the search terms and count them
						COUNT(*) AS tag_match_count
				FROM 
						productTags pt
				JOIN
						search_terms st ON pt.tagName ILIKE '%' || st.term || '%'
				GROUP BY 
						pt.ProductId
		)

		SELECT 
				COUNT(*) OVER() AS total,
				p.id, 
				p.name, 
				p.price, 
				p.descript, 
				p.quantity, 
				p.image,
				STRING_AGG(pt.tagName, ',') AS tags,
				
				-- Match counter to rank results with more matches higher
				(
						-- Count matches in name, description, and tags (from tag_matches)
						CASE 
								WHEN p.name ILIKE ANY (SELECT '%' || term || '%' FROM search_terms) THEN 1 
								ELSE 0 
						END +
						CASE 
								WHEN p.descript ILIKE ANY (SELECT '%' || term || '%' FROM search_terms) THEN 1 
								ELSE 0 
						END +
						COALESCE(tm.tag_match_count, 0) -- Use tag match count from subquery
				) AS match_count

		FROM 
				products p
		LEFT JOIN 
				productTags pt ON p.id = pt.ProductId
		LEFT JOIN
				tag_matches tm ON p.id = tm.ProductId
		WHERE 
				(
						p.name ILIKE ANY (SELECT '%' || term || '%' FROM search_terms) OR
						p.descript ILIKE ANY (SELECT '%' || term || '%' FROM search_terms) OR
						pt.tagName ILIKE ANY (SELECT '%' || term || '%' FROM search_terms)
				)
		GROUP BY 
				p.id, tm.tag_match_count
		ORDER BY 
				match_count DESC, -- Sort by match count first
				p.price           -- Then order by price
		LIMIT 10
		OFFSET $2;

		`
		rows, err = db.Query(query, term, offset)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		for rows.Next() {
			err := rows.Scan(&total, &id, &name, &price, &desc, &quantity, &imgByte, &tagsStr,&rank)
			if err != nil {
				log.Fatal(err)
			}

			imgStr := base64.StdEncoding.EncodeToString(imgByte)

			var tags []string
			if tagsStr != "" {
				tags = strings.Split(tagsStr, ",")
			}

			ProductList = append(ProductList, wc.Product{
				Id:       id,
				Name:     name,
				Price:    price,
				Desc:     desc,
				Quantity: quantity,
				Image:    imgStr,
				Tags:     tags,
			})
		}

		if err := rows.Err(); err != nil {
			log.Fatal(err)
		}

		return ProductList, total

  }else{
		for rows.Next() {
			err := rows.Scan(&total, &id, &name, &price, &desc, &quantity, &imgByte, &tagsStr, &rank)
			if err != nil {
				log.Fatal(err)
			}

			imgStr := base64.StdEncoding.EncodeToString(imgByte)

			var tags []string
			if tagsStr != "" {
				tags = strings.Split(tagsStr, ",")
			}

			ProductList = append(ProductList, wc.Product{
				Id:       id,
				Name:     name,
				Price:    price,
				Desc:     desc,
				Quantity: quantity,
				Image:    imgStr,
				Tags:     tags,
			})
		}

		if err := rows.Err(); err != nil {
			log.Fatal(err)
		}

		return ProductList, total
	}
}




var productTags = [60]string{
  "electronics", "clothing", "home-appliances", "books", "furniture",
  "sports", "toys", "automotive", "new-arrival", "best-seller",
  "limited-edition", "exclusive", "eco-friendly", "handmade", "premium",
  "affordable", "summer-sale", "winter-collection", "black-friday",
  "holiday-special", "back-to-school", "sale", "clearance", "discounted",
  "buy-one-get-one", "free-shipping", "flash-deal", "new", "refurbished",
  "used", "vintage", "open-box", "kids", "adults", "women", "men", "seniors",
  "unisex", "nike", "samsung", "apple", "sony", "adidas", "puma", "lego",
  "durable", "lightweight", "waterproof", "ergonomic", "energy-efficient",
  "fast-charging", "multi-purpose", "high-performance", "red", "blue", "black",
  "white", "modern", "minimalistic", "classic",
}
var productTagsMap = map[string]bool{
  "electronics": true, "clothing": true, "home-appliances": true, "books": true, "furniture": true,
  "sports": true, "toys": true, "automotive": true, "new-arrival": true, "best-seller": true,
  "limited-edition": true, "exclusive": true, "eco-friendly": true, "handmade": true, "premium": true,
  "affordable": true, "summer-sale": true, "winter-collection": true, "black-friday": true,
  "holiday-special": true, "back-to-school": true, "sale": true, "clearance": true, "discounted": true,
  "buy-one-get-one": true, "free-shipping": true, "flash-deal": true, "new": true, "refurbished": true,
  "used": true, "vintage": true, "open-box": true, "kids": true, "adults": true, "women": true,
  "men": true, "seniors": true, "unisex": true, "nike": true, "samsung": true, "apple": true,
  "sony": true, "adidas": true, "puma": true, "lego": true, "durable": true, "lightweight": true,
  "waterproof": true, "ergonomic": true, "energy-efficient": true, "fast-charging": true,
  "multi-purpose": true, "high-performance": true, "red": true, "blue": true, "black": true,
  "white": true, "modern": true, "minimalistic": true, "classic": true,
}


