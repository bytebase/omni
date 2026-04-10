use role securityadmin;

create or replace row access policy rap_sales_manager_regions_1 as (sales_region varchar) returns boolean ->
  'sales_executive_role' = current_role()
      or exists (
            select 1 from salesmanagerregions
              where sales_manager = current_role()
                and region = sales_region
          )
;
