-- pin: default gvenzl/oracle-free:23.26.3-slim-faststart sha256:f5ff19033860d662c821cb04eb10483fa94f14f78eae252d054291ea07028093 23.26.3.0.0

== TestVersionFromServer default#1
> query
| SELECT version_full FROM product_component_version WHERE product LIKE 'Oracle%'
< result
< column "VERSION_FULL" "NCHAR" "string"
< row
: string
| 23.26.3.0.0
