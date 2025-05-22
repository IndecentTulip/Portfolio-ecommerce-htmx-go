ALTER TABLE products ADD COLUMN search tsvector;

CREATE OR REPLACE FUNCTION update_search_column() 
RETURNS TRIGGER AS $$
BEGIN
  NEW.search := 
    setweight(to_tsvector('simple', NEW.name), 'A') || 
    setweight(to_tsvector('english', NEW.descript), 'B') || 
    setweight(to_tsvector('simple', 
        COALESCE(
            (SELECT string_agg(pt.tagName, ' ') 
             FROM productTags pt 
             JOIN tags t ON pt.TagName = t.tagName 
             WHERE pt.ProductId = NEW.id), 
            ''
        )
    ), 'C') :: tsvector;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_update_search 
BEFORE INSERT OR UPDATE ON products
FOR EACH ROW
EXECUTE FUNCTION update_search_column();

CREATE INDEX idx_products_search ON products USING gin(search);

UPDATE products
SET search = 
    setweight(to_tsvector('simple', name), 'A') || 
    setweight(to_tsvector('english', descript), 'B') || 
    setweight(to_tsvector('simple', 
        COALESCE(
            (SELECT string_agg(pt.tagName, ' ') 
             FROM productTags pt 
             JOIN tags t ON pt.TagName = t.tagName 
             WHERE pt.ProductId = products.id), 
            ''
        )
    ), 'C') :: tsvector;


