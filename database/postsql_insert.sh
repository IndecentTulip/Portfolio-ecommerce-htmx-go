#!/bin/bash

# ======
DB_NAME="ecommerce"
DB_USER="postgres"
DB_HOST="localhost"
IMG_PATH="./test.png"
# ======

# Check & Insert Products

IMG_HEX=$(xxd -p "$IMG_PATH" | tr -d '\n')

COUNT=$(psql -U "$DB_USER" -d "$DB_NAME" -t -c "SELECT COUNT(*) FROM products;" | xargs)

if [ "$COUNT" -ge 300 ]; then
  echo "Products already inserted."
else
  echo "Inserting products..."
  for i in $(seq 0 299); do
    UUID=$(psql -U "$DB_USER" -d "$DB_NAME" -t -c "SELECT gen_random_uuid();" | xargs)
    NAME="product $i"
    PRICE=$((25 + i))
    DESCRIPT="Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat."
    [ "$i" -gt 50 ] && DESCRIPT="DEMO TEST DEMO TEST DEMO TEST DEMO TEST DEMO TEST"
    [ "$i" -gt 100 ] && DESCRIPT="It is not knowledge, but the act of learning, not possession but the act of getting there, which grants the greatest enjoyment"
    [ "$i" -gt 200 ] && DESCRIPT="A computer is like air conditioning - it becomes useless when you open Windows"

    psql -U "$DB_USER" -d "$DB_NAME" -c \
    "INSERT INTO products (id, name, price, descript, quantity, image) 
     VALUES ('$UUID', '$NAME', $PRICE, '$DESCRIPT', 3, decode('$IMG_HEX', 'hex'))"
  done
  echo "Products inserted."
fi

# Insert Default Tags

TAGS=(
  "electronics" "clothing" "home-appliances" "books" "furniture"
  "sports" "toys" "automotive" "new-arrival" "best-seller"
  "limited-edition" "exclusive" "eco-friendly" "handmade" "premium"
  "affordable" "summer-sale" "winter-collection" "black-friday"
  "holiday-special" "back-to-school" "sale" "clearance" "discounted"
  "buy-one-get-one" "free-shipping" "flash-deal" "new" "refurbished"
  "used" "vintage" "open-box" "kids" "adults" "women" "men" "seniors"
  "unisex" "nike" "samsung" "apple" "sony" "adidas" "puma" "lego"
  "durable" "lightweight" "waterproof" "ergonomic" "energy-efficient"
  "fast-charging" "multi-purpose" "high-performance" "red" "blue" "black"
  "white" "modern" "minimalistic" "classic"
)

echo "Inserting tags..."
for TAG in "${TAGS[@]}"; do
  psql -U "$DB_USER" -d "$DB_NAME" -c \
  "INSERT INTO tags (tagname) VALUES ('$TAG') ON CONFLICT DO NOTHING;"
done
echo "Tags inserted."

# Assign Tags Randomly to Products

echo "Assigning tags to products..."

PRODUCT_IDS=($(psql -U "$DB_USER" -d "$DB_NAME" -Atc "SELECT id FROM products;"))
TAG_COUNT=${#TAGS[@]}

for PID in "${PRODUCT_IDS[@]}"; do
  TAG1=${TAGS[$((RANDOM % TAG_COUNT))]}
  TAG2=${TAGS[$((RANDOM % TAG_COUNT))]}

  while [ "$TAG1" == "$TAG2" ]; do
    TAG2=${TAGS[$((RANDOM % TAG_COUNT))]}
  done

  psql -U "$DB_USER" -d "$DB_NAME" -c \
  "INSERT INTO productTags (ProductId, TagName) VALUES 
   ('$PID', '$TAG1'), ('$PID', '$TAG2') 
   ON CONFLICT DO NOTHING;"

  echo "Tagged product $PID with ($TAG1, $TAG2)"
done

echo "Tag assignment complete!"

