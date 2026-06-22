-- Storage-Attached Index (SAI) examples

-- Basic SAI index
CREATE INDEX ON products (name) USING 'sai';

-- Custom SAI index with explicit name
CREATE CUSTOM INDEX product_name_idx ON products (name) USING 'StorageAttachedIndex';

-- Vector index for similarity search
CREATE CUSTOM INDEX embedding_idx ON products (embedding) 
    USING 'StorageAttachedIndex' 
    WITH {'similarity_function': 'cosine'};

-- Collection indexes with SAI (entire collection)
CREATE INDEX ON mytable (mymap) USING 'sai';