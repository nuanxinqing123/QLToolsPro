function r(r){return r.split(".")[0].replace("T"," ").replace("Z"," ")}function t(r){if("string"==typeof r)return r;try{return JSON.stringify(r)}catch(t){return""}}function n(r){if("number"==typeof r)return r;try{return JSON.parse(r)}catch(t){return 0}}export{r as d,t as n,n as s};
