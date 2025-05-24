ALTER TABLE products ADD COLUMN search tsvector;

CREATE OR REPLACE FUNCTION update_search_column() 
RETURNS TRIGGER AS $$
BEGIN
  NEW.search := 
    setweight(to_tsvector('simple', NEW.name || ' ' || get_prefixes_string(NEW.name) ), 'A') || 
    setweight(to_tsvector('english', NEW.descript), 'B') || 
    setweight(to_tsvector('simple', 
        COALESCE(
            (SELECT string_agg(pt.tagName || ' ' || get_prefixes_string(pt.tagName), ' ' ) 
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

CREATE OR REPLACE FUNCTION get_prefixes_string(word TEXT)
RETURNS TEXT AS $$
DECLARE
    result TEXT := '';
    i INTEGER;
BEGIN
    IF word IS NULL THEN
        RETURN '';
    END IF;

    FOR i IN 1..char_length(word) LOOP
        result := result || substring(word FROM 1 FOR i) || ' ';
    END LOOP;

    RETURN trim(result);
END;
$$ LANGUAGE plpgsql;


UPDATE products
SET search = 'update';


