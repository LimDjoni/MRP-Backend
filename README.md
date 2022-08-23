# Deli-AJE-Backend

# API DOCUMENTATION DELI AJE

URL = {{url}}/api/v1

## USER API

**REGISTER**
REGISTER_URL = URL/register
BODY =
- username = STRING , UNIQUE
- email = EMAIL , STRING , UNIQUE
- password = STRING, min 6, max 12

BODY_EXAMPLE =
```
{
    "username": "Hendri",
    "email": "hendri@gmail.com",
    "password": "halloinipassword"
}
```

RESPONSE =
```
{
    "ID": 4,
    "CreatedAt": "2022-08-22T11:43:52.374883+07:00",
    "UpdatedAt": "2022-08-22T11:43:52.374883+07:00",
    "DeletedAt": null,
    "username": "Hendri",
    "password": "$2a$04$JRPoCyuf4A5SVqFpqhR1nuJUN759CLqjFxwf.tE3ZWDa6ZOOPl.UG",
    "email": "hendrii@gmail.com"
  }
```
RESPONSE_ERROR_DUPLICATE =
```
{
    "error": "ERROR: duplicate key value violates unique constraint \"idx_users_email\" (SQLSTATE 23505)"
}
{
    "error": "ERROR: duplicate key value violates unique constraint \"idx_users_username\" (SQLSTATE 23505)"
}
```

RESPONSE_ERROR_VALIDATION =
```
{
    "errors": [
        {
            "FailedField": "RegisterUserInput.Username",
            "Tag": "required",
            "Value": ""
        },
        {
            "FailedField": "RegisterUserInput.Email",
            "Tag": "required",
            "Value": ""
        },
        {
            "FailedField": "RegisterUserInput.Password",
            "Tag": "required",
            "Value": ""
        }
    ]
}
```

RESPONSE_ERROR_VALIDATION_PASSWORD =
```
{
    "errors": [
        {
            "FailedField": "RegisterUserInput.Password",
            "Tag": "min",
            "Value": "123"
        },
        {
            "FailedField": "RegisterUserInput.Password",
            "Tag": "required",
            "Value": ""
        },
        {
            "FailedField": "RegisterUserInput.Password",
            "Tag": "max",
            "Value": "123123123123123123"
        }
    ]
}
```

RESPONSE_ERROR_VALIDATION_EMAIL =
```
{
    "errors": [
        {
            "FailedField": "RegisterUserInput.Email",
            "Tag": "required",
            "Value": ""
        },
        {
            "FailedField": "RegisterUserInput.Email",
            "Tag": "email",
            "Value": "adwadw"
        }
    ]
}
```

**LOGIN**
LOGIN_URL= URL/login
BODY =
- data = STRING -> email  / username used in register
- password = STRING

BODY_EXAMPLE =
```
{
    "data": "Hendri",
    "password": "halloinipassword"  
}
```

RESPONSE =
```
{
    "username": "Hendri",
    "email": "hendraiiii@gmail.com",
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImhlbmRyYWlpaWlAZ21haWwuY29tIiwiaWQiOjEwLCJ1c2VybmFtZSI6IkhlbmRyaSJ9.SbFwfQY_h6VCGN2HOpFFXEwIQKMe59ThpuDngbfAYqY"
}
```

RESPONSE_ERROR =
```
{
    "error": "wrong email / username / password"
}
```

# TRANSACTION API
TRANSACTION_URL = URL/transaction

**LIST DATA TRANSACTION DN**
URL = TRANSACTION_URL/list/dn
PARAMS =
- page number
- field string
- sort ASC / DESC
- ship_name string
- barge_name string
- shipping_from date
- shipping_to date
- quantity number

AUTH = BEARER ACCESS_TOKEN
EXAMPLE = {{TRANSACTION_URL}}/list/dn?page=1&sort=desc&ship_name=Black&barge_name=Black&shipping_to=2022-08-20&field=created_at&quantity=10.9&shipping_from=2022-08-10
RESPONSE =
```
{
      "limit": 10,
      "page": 1,
      "total_rows": 1,
      "total_pages": 1,
      "data": [
                {
                  "ID": 39,
                  "CreatedAt": "2022-08-19T15:12:11.412577+07:00",
                  "UpdatedAt": "2022-08-19T15:12:11.412577+07:00",
                  "DeletedAt": null,
                  "dmo_id": null,
                  "dmo": null,
                  "id_number": "DN-2022-8-0028",
                  "transaction_type": "DN",
                  "number": 28,
                  "shipping_date": "2022-08-12T00:00:00Z",
                  "quantity": 1023.122,
                  "ship_name": "AJE",
                  "barge_name": "",
                  "vessel_name": "",
                  "seller": "",
                  "customer_name": "",
                  "loading_port_name": "",
                  "loading_port_location": "",
                  "unloading_port_name": "",
                  "unloading_port_location": "",
                  "dmo_destination_port": "",
                  "skb_date": "2021-08-12T00:00:00Z",
                  "skb_number": "",
                  "skab_date": null,
                  "skab_number": "",
                  "bill_of_lading_date": null,
                  "bill_of_lading_number": "",
                  "royalty_rate": 0,
                  "dp_royalty_currency": "",
                  "dp_royalty_price": 0,
                  "dp_royalty_date": null,
                  "dp_royalty_ntpn": "",
                  "dp_royalty_billing_code": "",
                  "dp_royalty_total": 0,
                  "payment_dp_royalty_currency": "",
                  "payment_dp_royalty_price": 0,
                  "payment_dp_royalty_date": null,
                  "payment_dp_royalty_ntpn": "",
                  "payment_dp_royalty_billing_code": "",
                  "payment_dp_royalty_total": 0,
                  "lhv_date": null,
                  "lhv_number": "",
                  "surveyor_name": "",
                  "cow_date": null,
                  "cow_number": "",
                  "coa_date": null,
                  "coa_number": "",
                  "quality_tm_ar": 0,
                  "quality_im_adb": 0,
                  "quality_ash_ar": 0,
                  "quality_ash_adb": 0,
                  "quality_vm_adb": 0,
                  "quality_fc_adb": 0,
                  "quality_ts_ar": 0,
                  "quality_ts_adb": 0,
                  "quality_calories_ar": 0,
                  "quality_calories_adb": 0,
                  "barging_distance": 0,
                  "sales_system": "",
                  "invoice_date": null,
                  "invoice_number": "",
                  "invoice_price_unit": 0,
                  "invoice_price_total": 0,
                  "dmo_reconciliation_letter": "",
                  "contract_date": null,
                  "contract_number": "",
                  "dmo_buyer_name": "",
                  "dmo_industry_type": "",
                  "dmo_status_reconciliation_letter": "",
                  "in_a_trade_number": "",
                  "description_of_good": "",
                  "tarif_pos_hs": "",
                  "volume": 0,
                  "unit": "",
                  "value": 0,
                  "currency": "",
                  "in_a_trade_loading_port": "",
                  "destination_country": "",
                  "data_ls_export_date": null,
                  "data_ls_export_number": "",
                  "data_ska_coo_date": null,
                  "data_ska_coo_number": "",
                  "peb_number": "",
                  "peb_date": null,
                  "aju_number": "",
                  "dwt_value": 0,
                  "insurance_company_name": "",
                  "insurance_polis_number": "",
                  "navy_ship_name": "",
                  "navy_company_name": "",
                  "navy_imo_number": "",
                  "skb_document": "",
                  "skab_document": "",
                  "bl_document": "",
                  "royalti_provision_document": "",
                  "royalti_final_document": "",
                  "cow_document": "",
                  "coa_document": "",
                  "invoice_and_contract_document": "",
                  "lhv_document": ""
                }
              ]
}
```

**DETAIL TRANSACTION DN**
URL = TRANSACTION_URL/detail/dn/:id
AUTH = BEARER ACCESS_TOKEN
EXAMPLE = {{TRANSACTION_URL}}/detail/dn/5
RESPONSE =
```
{
    "ID": 5,
    "CreatedAt": "2022-08-15T11:42:48.724359+07:00",
    "UpdatedAt": "2022-08-19T11:45:07.580801+07:00",
    "DeletedAt": null,
    "dmo_id": null,
    "dmo": null,
    "id_number": "DN-2022-8-4",
    "transaction_type": "DN",
    "number": 4,
    "shipping_date": null,
    "quantity": 0,
    "ship_name": "MV BELACKPINKU",
    "barge_name": "BG BELAKPINK",
    "vessel_name": "",
    "seller": "",
    "customer_name": "AJE CUST",
    "loading_port_name": "",
    "loading_port_location": "",
    "unloading_port_name": "",
    "unloading_port_location": "",
    "dmo_destination_port": "",
    "skb_date": null,
    "skb_number": "",
    "skab_date": null,
    "skab_number": "",
    "bill_of_lading_date": null,
    "bill_of_lading_number": "",
    "royalty_rate": 0,
    "dp_royalty_currency": "",
    "dp_royalty_price": 0,
    "dp_royalty_date": null,
    "dp_royalty_ntpn": "",
    "dp_royalty_billing_code": "",
    "dp_royalty_total": 0,
    "payment_dp_royalty_currency": "",
    "payment_dp_royalty_price": 0,
    "payment_dp_royalty_date": null,
    "payment_dp_royalty_ntpn": "",
    "payment_dp_royalty_billing_code": "",
    "payment_dp_royalty_total": 0,
    "lhv_date": null,
    "lhv_number": "",
    "surveyor_name": "",
    "cow_date": null,
    "cow_number": "",
    "coa_date": null,
    "coa_number": "",
    "quality_tm_ar": 0,
    "quality_im_adb": 0,
    "quality_ash_ar": 0,
    "quality_ash_adb": 0,
    "quality_vm_adb": 0,
    "quality_fc_adb": 0,
    "quality_ts_ar": 0,
    "quality_ts_adb": 0,
    "quality_calories_ar": 0,
    "quality_calories_adb": 0,
    "barging_distance": 0,
    "sales_system": "",
    "invoice_date": null,
    "invoice_number": "12352SDGE21",
    "invoice_price_unit": 0,
    "invoice_price_total": 0,
    "dmo_reconciliation_letter": "",
    "contract_date": null,
    "contract_number": "",
    "dmo_buyer_name": "",
    "dmo_industry_type": "",
    "dmo_status_reconciliation_letter": "",
    "in_a_trade_number": "",
    "description_of_good": "",
    "tarif_pos_hs": "",
    "volume": 0,
    "unit": "",
    "value": 0,
    "currency": "",
    "in_a_trade_loading_port": "",
    "destination_country": "",
    "data_ls_export_date": null,
    "data_ls_export_number": "",
    "data_ska_coo_date": null,
    "data_ska_coo_number": "",
    "peb_number": "",
    "peb_date": null,
    "aju_number": "",
    "dwt_value": 0,
    "insurance_company_name": "",
    "insurance_polis_number": "",
    "navy_ship_name": "",
    "navy_company_name": "",
    "navy_imo_number": "",
    "skb_document": "https://deli-aje.s3.ap-southeast-1.amazonaws.com/DN-2022-8-4/skb.pdf",
    "skab_document": "",
    "bl_document": "",
    "royalti_provision_document": "",
    "royalti_final_document": "",
    "cow_document": "",
    "coa_document": "",
    "invoice_and_contract_document": "",
    "lhv_document": ""
}
```

RESPONSE_ERROR =
```
{
    "error": "record not found"
}
```

**CREATE TRANSACTION DN**
URL = TRANSACTION_URL/create/dn
AUTH = BEARER ACCESS_TOKEN
RULES =
- only define when variable has value otherwise data will automatically define not have value
- if there is no value in number data will automatically be 0
- if there is no value in string data will automatically be ""
- if there is no value in date data will automatically be null
- if Date is available -> Format YYYY/MM/DD

BODY =
```
{
    "coa_date": null,
    "cow_date": null,
    "lhv_date": null,
    "quantity": 1023.122,
    "skb_date": "2022-02-01",
    "ship_name": "AJE",
    "skab_date": null,
    "barge_name": "",
    "coa_number": "",
    "cow_number": "",
    "lhv_number": "",
    "skb_number": "",
    "skab_number": "",
    "vessel_name": "",
    "invoice_date": null,
    "royalty_rate": 0,
    "sales_system": "",
    "contract_date": null,
    "customer_name": "",
    "quality_tm_ar": 0,
    "quality_ts_ar": 0,
    "shipping_date": null,
    "surveyor_name": "",
    "dmo_buyer_name": "",
    "invoice_number": "",
    "quality_ash_ar": 0,
    "quality_fc_adb": 0,
    "quality_im_adb": 0,
    "quality_ts_adb": 0,
    "quality_vm_adb": 0,
    "contract_number": "",
    "dp_royalty_date": null,
    "dp_royalty_ntpn": "",
    "quality_ash_adb": 0,
    "barging_distance": 0,
    "dp_royalty_price": 0,
    "dp_royalty_total": 0,
    "dmo_industry_type": "",
    "loading_port_name": "",
    "invoice_price_unit": 0,
    "bill_of_lading_date": null,
    "dp_royalty_currency": "",
    "invoice_price_total": 0,
    "quality_calories_ar": 0,
    "unloading_port_name": "",
    "dmo_destination_port": "",
    "quality_calories_adb": 0,
    "bill_of_lading_number": "",
    "loading_port_location": "",
    "dp_royalty_billing_code": "",
    "payment_dp_royalty_date": null,
    "payment_dp_royalty_ntpn": "",
    "unloading_port_location": "",
    "payment_dp_royalty_price": 0,
    "payment_dp_royalty_total": 0,
    "dmo_reconciliation_letter": "",
    "payment_dp_royalty_currency": "",
    "payment_dp_royalty_billing_code": "",
    "dmo_status_reconciliation_letter": ""
}
```

RESPONSE =
```
{
    "ID": 45,
    "CreatedAt": "2022-08-22T14:08:11.484637+07:00",
    "UpdatedAt": "2022-08-22T14:08:11.484637+07:00",
    "DeletedAt": null,
    "dmo_id": null,
    "dmo": null,
    "id_number": "DN-2022-8-0034",
    "transaction_type": "DN",
    "number": 34,
    "shipping_date": null,
    "quantity": 1023.122,
    "ship_name": "AJE",
    "barge_name": "",
    "vessel_name": "",
    "seller": "",
    "customer_name": "",
    "loading_port_name": "",
    "loading_port_location": "",
    "unloading_port_name": "",
    "unloading_port_location": "",
    "dmo_destination_port": "",
    "skb_date": "2022-02-01",
    "skb_number": "",
    "skab_date": null,
    "skab_number": "",
    "bill_of_lading_date": null,
    "bill_of_lading_number": "",
    "royalty_rate": 0,
    "dp_royalty_currency": "",
    "dp_royalty_price": 0,
    "dp_royalty_date": null,
    "dp_royalty_ntpn": "",
    "dp_royalty_billing_code": "",
    "dp_royalty_total": 0,
    "payment_dp_royalty_currency": "",
    "payment_dp_royalty_price": 0,
    "payment_dp_royalty_date": null,
    "payment_dp_royalty_ntpn": "",
    "payment_dp_royalty_billing_code": "",
    "payment_dp_royalty_total": 0,
    "lhv_date": null,
    "lhv_number": "",
    "surveyor_name": "",
    "cow_date": null,
    "cow_number": "",
    "coa_date": null,
    "coa_number": "",
    "quality_tm_ar": 0,
    "quality_im_adb": 0,
    "quality_ash_ar": 0,
    "quality_ash_adb": 0,
    "quality_vm_adb": 0,
    "quality_fc_adb": 0,
    "quality_ts_ar": 0,
    "quality_ts_adb": 0,
    "quality_calories_ar": 0,
    "quality_calories_adb": 0,
    "barging_distance": 0,
    "sales_system": "",
    "invoice_date": null,
    "invoice_number": "",
    "invoice_price_unit": 0,
    "invoice_price_total": 0,
    "dmo_reconciliation_letter": "",
    "contract_date": null,
    "contract_number": "",
    "dmo_buyer_name": "",
    "dmo_industry_type": "",
    "dmo_status_reconciliation_letter": "",
    "in_a_trade_number": "",
    "description_of_good": "",
    "tarif_pos_hs": "",
    "volume": 0,
    "unit": "",
    "value": 0,
    "currency": "",
    "in_a_trade_loading_port": "",
    "destination_country": "",
    "data_ls_export_date": null,
    "data_ls_export_number": "",
    "data_ska_coo_date": null,
    "data_ska_coo_number": "",
    "peb_number": "",
    "peb_date": null,
    "aju_number": "",
    "dwt_value": 0,
    "insurance_company_name": "",
    "insurance_polis_number": "",
    "navy_ship_name": "",
    "navy_company_name": "",
    "navy_imo_number": "",
    "skb_document": "",
    "skab_document": "",
    "bl_document": "",
    "royalti_provision_document": "",
    "royalti_final_document": "",
    "cow_document": "",
    "coa_document": "",
    "invoice_and_contract_document": "",
    "lhv_document": ""
}
```

RESPONSE_ERROR =
```
{
    "errors": [
      {
          "FailedField": "DataTransactionInput.ShippingDate",
          "Tag": "DateValidation",
          "Value": "2022-23-01"
      },
      {
          "FailedField": "DataTransactionInput.SkbDate",
          "Tag": "DateValidation",
          "Value": "2022-23-01"
      },
      {
          "FailedField": "DataTransactionInput.SkabDate",
          "Tag": "DateValidation",
          "Value": "2022-23-01"
      },
      {
          "FailedField": "DataTransactionInput.BillOfLadingDate",
          "Tag": "DateValidation",
          "Value": "2022-23-01"
      },
      {
          "FailedField": "DataTransactionInput.DpRoyaltyDate",
          "Tag": "DateValidation",
          "Value": "2022-23-01"
      },
      {
          "FailedField": "DataTransactionInput.PaymentDpRoyaltyDate",
          "Tag": "DateValidation",
          "Value": "2022-23-01"
      },
      {
          "FailedField": "DataTransactionInput.LhvDate",
          "Tag": "DateValidation",
          "Value": "2022-23-01"
      },
      {
          "FailedField": "DataTransactionInput.CowDate",
          "Tag": "DateValidation",
          "Value": "2022-23-01"
      },
      {
          "FailedField": "DataTransactionInput.CoaDate",
          "Tag": "DateValidation",
          "Value": "2022-23-01"
      },
      {
          "FailedField": "DataTransactionInput.InvoiceDate",
          "Tag": "DateValidation",
          "Value": "2022-23-01"
      },
      {
          "FailedField": "DataTransactionInput.ContractDate",
          "Tag": "DateValidation",
          "Value": "2022-23-01"
      }
    ]
  }
```

**DELETE TRANSACTION DN**
URL = TRANSACTION_URL/delete/dn/:id
AUTH = BEARER ACCESS_TOKEN
EXAMPLE = TRANSACTION_URL/delete/dn/5
RESPONSE =
```
{
    "message": "success delete transaction"
}
```
RESPONSE_ERROR =
```
{
    "error": "record not found",
    "message": "failed to delete transaction"
}
```

**UPDATE DATA TRANSACTION DN**
URL = TRANSACTION_URL/update/dn/:id
AUTH = BEARER ACCESS_TOKEN
RULES =
- can copy response of detail transaction and edit value if want update
- only data needed to create can be updated (not document)

BODY =
```
{
    "coa_date": null,
    "cow_date": null,
    "lhv_date": null,
    "quantity": 1023.122,
    "skb_date": "2022-02-01",
    "ship_name": "AJE",
    "skab_date": null,
    "barge_name": "",
    "coa_number": "",
    "cow_number": "",
    "lhv_number": "",
    "skb_number": "",
    "skab_number": "",
    "vessel_name": "",
    "invoice_date": null,
    "royalty_rate": 0,
    "sales_system": "",
    "contract_date": null,
    "customer_name": "",
    "quality_tm_ar": 0,
    "quality_ts_ar": 0,
    "shipping_date": null,
    "surveyor_name": "",
    "dmo_buyer_name": "",
    "invoice_number": "",
    "quality_ash_ar": 0,
    "quality_fc_adb": 0,
    "quality_im_adb": 0,
    "quality_ts_adb": 0,
    "quality_vm_adb": 0,
    "contract_number": "",
    "dp_royalty_date": null,
    "dp_royalty_ntpn": "",
    "quality_ash_adb": 0,
    "barging_distance": 0,
    "dp_royalty_price": 0,
    "dp_royalty_total": 0,
    "dmo_industry_type": "",
    "loading_port_name": "",
    "invoice_price_unit": 0,
    "bill_of_lading_date": null,
    "dp_royalty_currency": "",
    "invoice_price_total": 0,
    "quality_calories_ar": 0,
    "unloading_port_name": "",
    "dmo_destination_port": "",
    "quality_calories_adb": 0,
    "bill_of_lading_number": "",
    "loading_port_location": "",
    "dp_royalty_billing_code": "",
    "payment_dp_royalty_date": null,
    "payment_dp_royalty_ntpn": "",
    "unloading_port_location": "",
    "payment_dp_royalty_price": 0,
    "payment_dp_royalty_total": 0,
    "dmo_reconciliation_letter": "",
    "payment_dp_royalty_currency": "",
    "payment_dp_royalty_billing_code": "",
    "dmo_status_reconciliation_letter": ""
}
```

BODY_EXAMPLE =
```
{
    "ID": 6,
    "CreatedAt": "2022-08-15T11:43:11.031048+07:00",
    "UpdatedAt": "2022-08-15T11:43:11.031048+07:00",
    "DeletedAt": null,
    "dmo_id": null,
    "dmo": null,
    "id_number": "DN-2022-8-5",
    "transaction_type": "DN",
    "number": 5,
    "shipping_date": null,
    "quantity": 0,
    "ship_name": "BG GARUDA",
    "barge_name": "",
    "vessel_name": "",
    "seller": "",
    "customer_name": "",
    "loading_port_name": "",
    "loading_port_location": "",
    "unloading_port_name": "",
    "unloading_port_location": "",
    "dmo_destination_port": "",
    "skb_date": null,
    "skb_number": "",
    "skab_date": null,
    "skab_number": "",
    "bill_of_lading_date": null,
    "bill_of_lading_number": "",
    "royalty_rate": 0,
    "dp_royalty_currency": "",
    "dp_royalty_price": 0,
    "dp_royalty_date": null,
    "dp_royalty_ntpn": "",
    "dp_royalty_billing_code": "",
    "dp_royalty_total": 0,
    "payment_dp_royalty_currency": "",
    "payment_dp_royalty_price": 0,
    "payment_dp_royalty_date": null,
    "payment_dp_royalty_ntpn": "",
    "payment_dp_royalty_billing_code": "",
    "payment_dp_royalty_total": 0,
    "lhv_date": null,
    "lhv_number": "",
    "surveyor_name": "",
    "cow_date": null,
    "cow_number": "",
    "coa_date": null,
    "coa_number": "",
    "quality_tm_ar": 0,
    "quality_im_adb": 0,
    "quality_ash_ar": 0,
    "quality_ash_adb": 0,
    "quality_vm_adb": 0,
    "quality_fc_adb": 0,
    "quality_ts_ar": 0,
    "quality_ts_adb": 0,
    "quality_calories_ar": 0,
    "quality_calories_adb": 0,
    "barging_distance": 0,
    "sales_system": "",
    "invoice_date": null,
    "invoice_number": "",
    "invoice_price_unit": 0,
    "invoice_price_total": 0,
    "dmo_reconciliation_letter": "",
    "contract_date": null,
    "contract_number": "",
    "dmo_buyer_name": "",
    "dmo_industry_type": "",
    "dmo_status_reconciliation_letter": "",
    "in_a_trade_number": "",
    "description_of_good": "",
    "tarif_pos_hs": "",
    "volume": 0,
    "unit": "",
    "value": 0,
    "currency": "",
    "in_a_trade_loading_port": "",
    "destination_country": "",
    "data_ls_export_date": null,
    "data_ls_export_number": "",
    "data_ska_coo_date": null,
    "data_ska_coo_number": "",
    "peb_number": "",
    "peb_date": null,
    "aju_number": "",
    "dwt_value": 0,
    "insurance_company_name": "",
    "insurance_polis_number": "",
    "navy_ship_name": "",
    "navy_company_name": "",
    "navy_imo_number": "",
    "skb_document": "",
    "skab_document": "",
    "bl_document": "",
    "royalti_provision_document": "",
    "royalti_final_document": "",
    "cow_document": "",
    "coa_document": "",
    "invoice_and_contract_document": "",
    "lhv_document": ""
}
```

RESPONSE =
```
{
    "ID": 6,
    "CreatedAt": "2022-08-15T11:43:11.031048+07:00",
    "UpdatedAt": "2022-08-22T14:21:55.799549+07:00",
    "DeletedAt": null,
    "dmo_id": null,
    "dmo": null,
    "id_number": "DN-2022-8-5",
    "transaction_type": "DN",
    "number": 5,
    "shipping_date": null,
    "quantity": 0,
    "ship_name": "BG GARUDA",
    "barge_name": "",
    "vessel_name": "",
    "seller": "",
    "customer_name": "",
    "loading_port_name": "",
    "loading_port_location": "",
    "unloading_port_name": "",
    "unloading_port_location": "",
    "dmo_destination_port": "",
    "skb_date": null,
    "skb_number": "",
    "skab_date": null,
    "skab_number": "",
    "bill_of_lading_date": null,
    "bill_of_lading_number": "",
    "royalty_rate": 0,
    "dp_royalty_currency": "",
    "dp_royalty_price": 0,
    "dp_royalty_date": null,
    "dp_royalty_ntpn": "",
    "dp_royalty_billing_code": "",
    "dp_royalty_total": 0,
    "payment_dp_royalty_currency": "",
    "payment_dp_royalty_price": 0,
    "payment_dp_royalty_date": null,
    "payment_dp_royalty_ntpn": "",
    "payment_dp_royalty_billing_code": "",
    "payment_dp_royalty_total": 0,
    "lhv_date": null,
    "lhv_number": "",
    "surveyor_name": "",
    "cow_date": null,
    "cow_number": "",
    "coa_date": null,
    "coa_number": "",
    "quality_tm_ar": 0,
    "quality_im_adb": 0,
    "quality_ash_ar": 0,
    "quality_ash_adb": 0,
    "quality_vm_adb": 0,
    "quality_fc_adb": 0,
    "quality_ts_ar": 0,
    "quality_ts_adb": 0,
    "quality_calories_ar": 0,
    "quality_calories_adb": 0,
    "barging_distance": 0,
    "sales_system": "",
    "invoice_date": null,
    "invoice_number": "",
    "invoice_price_unit": 0,
    "invoice_price_total": 0,
    "dmo_reconciliation_letter": "",
    "contract_date": null,
    "contract_number": "",
    "dmo_buyer_name": "",
    "dmo_industry_type": "",
    "dmo_status_reconciliation_letter": "",
    "in_a_trade_number": "",
    "description_of_good": "",
    "tarif_pos_hs": "",
    "volume": 0,
    "unit": "",
    "value": 0,
    "currency": "",
    "in_a_trade_loading_port": "",
    "destination_country": "",
    "data_ls_export_date": null,
    "data_ls_export_number": "",
    "data_ska_coo_date": null,
    "data_ska_coo_number": "",
    "peb_number": "",
    "peb_date": null,
    "aju_number": "",
    "dwt_value": 0,
    "insurance_company_name": "",
    "insurance_polis_number": "",
    "navy_ship_name": "",
    "navy_company_name": "",
    "navy_imo_number": "",
    "skb_document": "",
    "skab_document": "",
    "bl_document": "",
    "royalti_provision_document": "",
    "royalti_final_document": "",
    "cow_document": "",
    "coa_document": "",
    "invoice_and_contract_document": "",
    "lhv_document": ""
}
```

RESPONSE_ERROR =
```
{
    "error": "record not found",
    "message": "failed to update transaction"
}
```

**UPDATE DOCUMENT TRANSACTION DN**
URL = TRANSACTION_URL/update/document/dn/:id/:type
AUTH = BEARER ACCESS_TOKEN
RULES =
- only type provided can be upload (skb, skab, bl, royalti_provision, royalti_final, cow, coa, invoice, lhv)

BODY =
```
{
    document -> file to upload (must be pdf)
}
```

RESPONSE =
```
{
    "ID": 37,
    "CreatedAt": "2022-08-19T09:22:16.099263+07:00",
    "UpdatedAt": "2022-08-22T14:24:55.632676+07:00",
    "DeletedAt": null,
    "dmo_id": null,
    "dmo": null,
    "id_number": "DN-2022-8-0026",
    "transaction_type": "DN",
    "number": 26,
    "shipping_date": null,
    "quantity": 1023.122,
    "ship_name": "AJE",
    "barge_name": "",
    "vessel_name": "",
    "seller": "",
    "customer_name": "",
    "loading_port_name": "",
    "loading_port_location": "",
    "unloading_port_name": "",
    "unloading_port_location": "",
    "dmo_destination_port": "",
    "skb_date": "2021-08-12T00:00:00Z",
    "skb_number": "",
    "skab_date": null,
    "skab_number": "",
    "bill_of_lading_date": null,
    "bill_of_lading_number": "",
    "royalty_rate": 0,
    "dp_royalty_currency": "",
    "dp_royalty_price": 0,
    "dp_royalty_date": null,
    "dp_royalty_ntpn": "",
    "dp_royalty_billing_code": "",
    "dp_royalty_total": 0,
    "payment_dp_royalty_currency": "",
    "payment_dp_royalty_price": 0,
    "payment_dp_royalty_date": null,
    "payment_dp_royalty_ntpn": "",
    "payment_dp_royalty_billing_code": "",
    "payment_dp_royalty_total": 0,
    "lhv_date": null,
    "lhv_number": "",
    "surveyor_name": "",
    "cow_date": null,
    "cow_number": "",
    "coa_date": null,
    "coa_number": "",
    "quality_tm_ar": 0,
    "quality_im_adb": 0,
    "quality_ash_ar": 0,
    "quality_ash_adb": 0,
    "quality_vm_adb": 0,
    "quality_fc_adb": 0,
    "quality_ts_ar": 0,
    "quality_ts_adb": 0,
    "quality_calories_ar": 0,
    "quality_calories_adb": 0,
    "barging_distance": 0,
    "sales_system": "",
    "invoice_date": null,
    "invoice_number": "",
    "invoice_price_unit": 0,
    "invoice_price_total": 0,
    "dmo_reconciliation_letter": "",
    "contract_date": null,
    "contract_number": "",
    "dmo_buyer_name": "",
    "dmo_industry_type": "",
    "dmo_status_reconciliation_letter": "",
    "in_a_trade_number": "",
    "description_of_good": "",
    "tarif_pos_hs": "",
    "volume": 0,
    "unit": "",
    "value": 0,
    "currency": "",
    "in_a_trade_loading_port": "",
    "destination_country": "",
    "data_ls_export_date": null,
    "data_ls_export_number": "",
    "data_ska_coo_date": null,
    "data_ska_coo_number": "",
    "peb_number": "",
    "peb_date": null,
    "aju_number": "",
    "dwt_value": 0,
    "insurance_company_name": "",
    "insurance_polis_number": "",
    "navy_ship_name": "",
    "navy_company_name": "",
    "navy_imo_number": "",
    "skb_document": "https://deli-aje.s3.ap-southeast-1.amazonaws.com/DN-2022-8-0026/skb.pdf",
    "skab_document": "https://deli-aje.s3.ap-southeast-1.amazonaws.com/DN-2022-8-0026/skab.pdf",
    "bl_document": "https://deli-aje.s3.ap-southeast-1.amazonaws.com/DN-2022-8-0026/bl.pdf",
    "royalti_provision_document": "https://deli-aje.s3.ap-southeast-1.amazonaws.com/DN-2022-8-0026/royalti_provision.pdf",
    "royalti_final_document": "https://deli-aje.s3.ap-southeast-1.amazonaws.com/DN-2022-8-0026/royalti_final.pdf",
    "cow_document": "https://deli-aje.s3.ap-southeast-1.amazonaws.com/DN-2022-8-0026/cow.pdf",
    "coa_document": "https://deli-aje.s3.ap-southeast-1.amazonaws.com/DN-2022-8-0026/coa.pdf",
    "invoice_and_contract_document": "https://deli-aje.s3.ap-southeast-1.amazonaws.com/DN-2022-8-0026/invoice.pdf",
    "lhv_document": "https://deli-aje.s3.ap-southeast-1.amazonaws.com/DN-2022-8-0026/lhv.pdf"
}
```

RESPONSE_DOCUMENT_TYPE_ERROR =
```
{
    "error": "document type not found"
    "message": "failed to upload document"
}
```

RESPONSE_ERROR =
```
{
    "error": "record not found",
    "message": "failed to upload document"
}
```

RESPOSE_ERROR_DOCUMENT_TYPE_NOT_PDF =
```
{
    "error": "document must be pdf",
    "message": "failed to upload document"
}
```
