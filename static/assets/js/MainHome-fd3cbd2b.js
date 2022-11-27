import{_ as e,a}from"./footer-a4f2d620.js";import{r as s}from"./index-5edb5878.js";import{c as t}from"./store-8e5b6e13.js";import{s as l}from"./pinia-1ca0c29b.js";import{a4 as i,a5 as o,a6 as n,a7 as r,a8 as p,a9 as c,aa as u,ab as m,ac as d,ad as f,ae as g,af as j}from"./@ant-design-477981d5.js";import{r as b,o as k,a4 as y,a6 as _,a3 as v,f as h,c as C,u as w,i as M,a5 as H,M as x,$ as S,_ as P,L as K}from"./@vue-11129043.js";import"./recharge-pop-8b070e01.js";import"./p-center-modal-84ed26df.js";import"./_plugin-vue_export-helper-9b9a8a5b.js";import"./recharge-data-b14aa3fe.js";import"./common-3c8652c5.js";import"./ant-design-vue-93fc40b2.js";import"./@babel-eab0ef53.js";import"./resize-observer-polyfill-9cd09a67.js";import"./@ctrl-16df70a4.js";import"./dom-align-6a5270eb.js";import"./lodash-es-0ceb8576.js";import"./dayjs-c306acb6.js";import"./async-validator-604317c1.js";import"./scroll-into-view-if-needed-8ce8502d.js";import"./compute-scroll-into-view-cce79123.js";import"./axios-93d3568f.js";import"./qs-8fb0a9f1.js";import"./side-channel-ee547e73.js";import"./get-intrinsic-53528089.js";import"./has-symbols-1f359e75.js";import"./function-bind-c930bb92.js";import"./has-03e8e28c.js";import"./call-bind-566c57e8.js";import"./object-inspect-5c6480f3.js";import"./vue-router-36397834.js";import"./vue3-colorpicker-e1559e09.js";import"./vue-types-0fd36d85.js";import"./is-plain-object-39134198.js";import"./tinycolor2-e232e212.js";import"./@vueuse-94329f85.js";import"./@aesoper-316802a3.js";import"./vue3-angle-2884cf46.js";import"./gradient-parser-c9367eab.js";import"./@popperjs-31624eb1.js";import"./vue-demi-a81ff0a7.js";const R={key:1,class:"page-logo text-center font-weight"},$=P("span",null,"首页",-1),z=P("span",null,"总览",-1),D=P("span",null,"面板管理",-1),L=P("span",null,"变量管理",-1),U=P("span",null,"消息推送管理",-1),q=P("span",null,"容器管理",-1),E=P("span",null,"插件管理",-1),I=P("span",null,"用户管理",-1),T=P("span",null,"用户变量管理",-1),V=P("span",null,"卡密管理",-1),A=K("数据查询"),B=K("充值数据"),F=K("上传记录"),G=P("span",null,"网站设置",-1),J={__name:"MainHome",setup(K){const{isMobile:J,isCollapsed:N}=l(t),O=b(t.pageKeys),Q=b(null);k((()=>{let e=Q.value.clientHeight||Q.value.$el.clientHeight;t.setRouterPageHeight(e),window.addEventListener("resize",(()=>{let e=Q.value.clientHeight||Q.value.$el.clientHeight;t.setRouterPageHeight(e)}))}));const W=(e,a)=>{J.value&&t.setCollapsed(!0),s.push({name:e})};return(s,l)=>{const b=v("a-image"),k=v("a-menu-item"),K=v("a-sub-menu"),J=v("a-menu"),X=v("a-layout-sider"),Y=v("router-view"),Z=v("a-layout-content"),ee=v("a-layout");return h(),y(ee,{style:{"min-height":"100vh"}},{default:_((()=>[C(X,{collapsed:w(N),"onUpdate:collapsed":l[14]||(l[14]=e=>x(N)?N.value=e:null),class:S(w(t).siteSettings.web_bg?"slide-hide-bg":""),collapsible:!0},{default:_((()=>[w(t).siteSettings.web_logo?(h(),y(b,{key:0,width:70,preview:!1,src:w(t).siteSettings.web_logo},null,8,["src"])):(h(),M("div",R,H(w(t).siteSettings.web_title||"青龙Tools Pro"),1)),C(J,{selectedKeys:O.value,"onUpdate:selectedKeys":l[13]||(l[13]=e=>O.value=e),theme:"dark",mode:"inline"},{default:_((()=>[C(k,{key:"13",onClick:l[0]||(l[0]=e=>W("adminHome"))},{default:_((()=>[C(w(i)),$])),_:1}),C(k,{key:"1",onClick:l[1]||(l[1]=e=>W("home"))},{default:_((()=>[C(w(o)),z])),_:1}),C(k,{key:"2",onClick:l[2]||(l[2]=e=>W("panelManagement"))},{default:_((()=>[C(w(n)),D])),_:1}),C(k,{key:"8",onClick:l[3]||(l[3]=e=>W("variableManagement"))},{default:_((()=>[C(w(r)),L])),_:1}),C(k,{key:"4",onClick:l[4]||(l[4]=e=>W("messagePushManagement"))},{default:_((()=>[C(w(p)),U])),_:1}),C(k,{key:"12",onClick:l[5]||(l[5]=e=>W("containerManagement"))},{default:_((()=>[C(w(c)),q])),_:1}),C(k,{key:"10",onClick:l[6]||(l[6]=e=>W("plugInManagement"))},{default:_((()=>[C(w(u)),E])),_:1}),C(k,{key:"3",onClick:l[7]||(l[7]=e=>W("userManagement"))},{default:_((()=>[C(w(m)),I])),_:1}),C(k,{key:"5",onClick:l[8]||(l[8]=e=>W("userVariableManagement"))},{default:_((()=>[C(w(d)),T])),_:1}),C(k,{key:"9",onClick:l[9]||(l[9]=e=>W("cardSecretManagement"))},{default:_((()=>[C(w(f)),V])),_:1}),C(K,{key:"sub1"},{icon:_((()=>[C(w(g))])),title:_((()=>[A])),default:_((()=>[C(k,{key:"6",onClick:l[10]||(l[10]=e=>W("rechargeData"))},{default:_((()=>[B])),_:1}),C(k,{key:"7",onClick:l[11]||(l[11]=e=>W("uploadData"))},{default:_((()=>[F])),_:1})])),_:1}),C(k,{key:"11",onClick:l[12]||(l[12]=e=>W("webSettings"))},{default:_((()=>[C(w(j)),G])),_:1})])),_:1},8,["selectedKeys"])])),_:1},8,["collapsed","class"]),C(ee,null,{default:_((()=>[C(e),C(Z,null,{default:_((()=>[P("div",{style:{height:"100%"},ref_key:"routerPageRef",ref:Q,class:"flex flex-column page-container"},[C(Y)],512)])),_:1}),C(a)])),_:1})])),_:1})}}};export{J as default};