import{N as e,k as a,w as t,f as l,a4 as s,a6 as i,c as o,a8 as r,u as n,L as p,a3 as u,r as d,i as m,F as v,J as c,a2 as f}from"./@vue-11129043.js";import{p as b}from"./p-center-modal-b39d0887.js";import"./dayjs-04e8f6f9.js";import{L as j,M as _,N as y}from"./index-e150f431.js";import{m as g,F as h,E as x}from"./ant-design-vue-05054bf0.js";import{_ as k}from"./page-container-fd59b876.js";import{d as C}from"./common-3c8652c5.js";import{S as U,ak as S}from"./@ant-design-477981d5.js";import"./pinia-1ca0c29b.js";import"./vue-demi-a81ff0a7.js";import"./store-bc016a64.js";import"./_plugin-vue_export-helper-9b9a8a5b.js";import"./@babel-eab0ef53.js";import"./axios-93d3568f.js";import"./qs-8fb0a9f1.js";import"./side-channel-ee547e73.js";import"./get-intrinsic-53528089.js";import"./has-symbols-1f359e75.js";import"./function-bind-c930bb92.js";import"./has-03e8e28c.js";import"./call-bind-566c57e8.js";import"./object-inspect-5c6480f3.js";import"./vue-router-36397834.js";import"./vue3-colorpicker-e1559e09.js";import"./vue-types-0fd36d85.js";import"./is-plain-object-39134198.js";import"./tinycolor2-e232e212.js";import"./@vueuse-94329f85.js";import"./@aesoper-316802a3.js";import"./vue3-angle-2884cf46.js";import"./gradient-parser-c9367eab.js";import"./lodash-es-0ceb8576.js";import"./@popperjs-31624eb1.js";import"./resize-observer-polyfill-9cd09a67.js";import"./@ctrl-16df70a4.js";import"./dom-align-6a5270eb.js";import"./async-validator-604317c1.js";import"./scroll-into-view-if-needed-8ce8502d.js";import"./compute-scroll-into-view-cce79123.js";const w=p("取消"),A=p("提交"),F={__name:"editOrAdd",props:{dataObj:Object,visible:Boolean},emits:["update:visible","updateData"],setup(p,{emit:d}){const m=p,{visible:v,dataObj:c}=e(m),f=()=>{d("update:visible",!1)},_=a({remarks:"",value:""}),y={value:[{required:!0,trigger:"change"}],remarks:[{required:!0,trigger:"change"}]},h=e=>{},x=(...e)=>{},k={labelCol:{span:7},wrapperCol:{span:15}},C=e=>{const a={..._,panel_name:c.value.panel_name,id:c.value.id,_id:c.value._id,name:c.value.name};j({data:a}).then((()=>{g.success("操作成功"),d("updateData",a),f()}))};t(v,((e,a,t)=>{v.value&&U()}));const U=()=>{_.remarks=c.value.remarks||"",_.value=c.value.value||""};return(e,a)=>{const t=u("a-input"),p=u("a-form-item"),d=u("a-divider"),m=u("a-button"),c=u("a-form");return l(),s(b,{modalVisible:n(v),isFooter:!1,onClose:f,title:"变量修改"},{content:i((()=>[o(c,r({ref:"formRef",name:"custom-validation",model:_,rules:y},k,{onValidate:x,onFinishFailed:h,onFinish:C}),{default:i((()=>[o(p,{label:"变量值",name:"value"},{default:i((()=>[o(t,{value:_.value,"onUpdate:value":a[0]||(a[0]=e=>_.value=e)},null,8,["value"])])),_:1}),o(p,{label:"变量备注",name:"remarks"},{default:i((()=>[o(t,{value:_.remarks,"onUpdate:value":a[1]||(a[1]=e=>_.remarks=e)},null,8,["value"])])),_:1}),o(d),o(p,{"wrapper-col":{span:12,offset:12},class:"timed-start-button"},{default:i((()=>[o(m,{onClick:f},{default:i((()=>[w])),_:1}),o(m,{type:"primary",style:{"margin-left":"20px"},"html-type":"submit"},{default:i((()=>[A])),_:1})])),_:1})])),_:1},16,["model"])])),_:1},8,["modalVisible"])}}},I=p("面板"),O=p("用户"),D=p(" 搜索 "),E=p("重置 "),z=p("修改 "),q=p(" 是否确认删除？ "),L=p("删除 "),M={__name:"index",setup(e){const t=d({}),s=d({}),r=d(!1),p=h.useForm,b=d(0),j=d(1),w=d(10);x.PRESENTED_IMAGE_SIMPLE;const A=a({s:"",type:"panel"}),M=d(!1),{resetFields:N,validate:V,validateInfos:K}=p(A),P=[{title:"面板名称",dataIndex:"panel_name",width:200},{title:"变量名称",dataIndex:"name",key:"name"},{title:"变量值",dataIndex:"value",key:"value"},{title:"操作",dataIndex:"operation",customKey:"operation"}],R=d([]),T=d([]),B=()=>{let e=j.value*w.value;e>b.value&&(e=b.value),R.value=T.value.slice((j.value-1)*w.value,e)},G=e=>{e&&(j.value=1);let a=s.value;a.page=j.value,a.quantity=w.value,_({data:A,splicingData:a}).then((e=>{T.value=e.filter((e=>e.panel_env&&e.panel_name)).map((e=>e.panel_env?e.panel_env.map((a=>(a.panel_name=e.panel_name,a.CreatedAt&&(a.CreatedAt=C(a.CreatedAt)),a.UpdatedAt&&(a.UpdatedAt=C(a.UpdatedAt)),a))):[])).flat(),b.value=T.value.length,B()}))};s.value=A;const J=()=>{N(),M.value=!1,G(!0)},H=()=>{V().then((e=>{s.value=A,G(!0)})).catch((e=>{}))};return(e,a)=>{const s=u("a-radio"),p=u("a-radio-group"),d=u("a-form-item"),_=u("a-input"),h=u("a-button"),x=u("a-form"),C=u("a-popconfirm");return l(),m(v,null,[o(F,{visible:r.value,"onUpdate:visible":a[0]||(a[0]=e=>r.value=e),onUpdateData:a[1]||(a[1]=e=>G(!0)),dataObj:t.value},null,8,["visible","dataObj"]),o(k,{columns:P,pageSize:w.value,"onUpdate:pageSize":a[4]||(a[4]=e=>w.value=e),current:j.value,"onUpdate:current":a[5]||(a[5]=e=>j.value=e),total:b.value,"onUpdate:total":a[6]||(a[6]=e=>b.value=e),isTable:!0,isSearch:!0,onOnShowSizeChange:B,onInitData:G,dataSource:R.value},{search:i((()=>[o(x,{class:"flex flex-warp",model:A},{default:i((()=>[o(d,{label:"搜索类型:",name:"type"},{default:i((()=>[o(p,{value:A.type,"onUpdate:value":a[2]||(a[2]=e=>A.type=e)},{default:i((()=>[o(s,{value:"panel",name:"type"},{default:i((()=>[I])),_:1}),o(s,{value:"user",name:"type"},{default:i((()=>[O])),_:1})])),_:1},8,["value"])])),_:1}),o(d,{label:"关键字:",name:"s"},{default:i((()=>[o(_,{value:A.s,"onUpdate:value":a[3]||(a[3]=e=>A.s=e),placeholder:"请输入关键字"},null,8,["value"])])),_:1}),o(d,null,{default:i((()=>[o(h,{type:"primary",onClick:c(H,["prevent"]),class:"filter-search"},{default:i((()=>[o(n(U)),D])),_:1},8,["onClick"]),o(h,{style:{"margin-left":"10px"},class:"filter-reset",onClick:J},{default:i((()=>[o(n(S)),E])),_:1})])),_:1})])),_:1},8,["model"])])),bodyCell:i((({text:e,record:a,index:s,column:n})=>["operation"===n.customKey?(l(),m(v,{key:0},[o(h,{type:"primary",onClick:c((e=>{return l=a,r.value=!0,void(t.value=l?{...l,title:"面板编辑"}:{title:"面板新增"});var l}),["stop"]),style:{"margin-left":"10px","margin-bottom":"10px"},shape:"round"},{default:i((()=>[z])),_:2},1032,["onClick"]),o(C,{placement:"topLeft","ok-text":"是","cancel-text":"否",onConfirm:e=>{var t;y({data:{panel_name:(t=a).panel_name,id:t.id,_id:t._id}}).then((()=>{g.success("操作成功!"),G(!0)}))}},{title:i((()=>[q])),default:i((()=>[o(h,{type:"danger",style:{"margin-left":"10px","margin-bottom":"10px"},shape:"round"},{default:i((()=>[L])),_:1})])),_:2},1032,["onConfirm"])],64)):f("",!0)])),_:1},8,["pageSize","current","total","dataSource"])],64)}}};export{M as default};
