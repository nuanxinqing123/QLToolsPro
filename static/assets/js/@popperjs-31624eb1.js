var e="top",t="bottom",n="right",r="left",o=[e,t,n,r],i=o.reduce((function(e,t){return e.concat([t+"-start",t+"-end"])}),[]),a=[].concat(o,["auto"]).reduce((function(e,t){return e.concat([t,t+"-start",t+"-end"])}),[]),s=["beforeRead","read","afterRead","beforeMain","main","afterMain","beforeWrite","write","afterWrite"];function f(e){return e?(e.nodeName||"").toLowerCase():null}function c(e){if(null==e)return window;if("[object Window]"!==e.toString()){var t=e.ownerDocument;return t&&t.defaultView||window}return e}function p(e){return e instanceof c(e).Element||e instanceof Element}function u(e){return e instanceof c(e).HTMLElement||e instanceof HTMLElement}function l(e){return"undefined"!=typeof ShadowRoot&&(e instanceof c(e).ShadowRoot||e instanceof ShadowRoot)}const d={name:"applyStyles",enabled:!0,phase:"write",fn:function(e){var t=e.state;Object.keys(t.elements).forEach((function(e){var n=t.styles[e]||{},r=t.attributes[e]||{},o=t.elements[e];u(o)&&f(o)&&(Object.assign(o.style,n),Object.keys(r).forEach((function(e){var t=r[e];!1===t?o.removeAttribute(e):o.setAttribute(e,!0===t?"":t)})))}))},effect:function(e){var t=e.state,n={popper:{position:t.options.strategy,left:"0",top:"0",margin:"0"},arrow:{position:"absolute"},reference:{}};return Object.assign(t.elements.popper.style,n.popper),t.styles=n,t.elements.arrow&&Object.assign(t.elements.arrow.style,n.arrow),function(){Object.keys(t.elements).forEach((function(e){var r=t.elements[e],o=t.attributes[e]||{},i=Object.keys(t.styles.hasOwnProperty(e)?t.styles[e]:n[e]).reduce((function(e,t){return e[t]="",e}),{});u(r)&&f(r)&&(Object.assign(r.style,i),Object.keys(o).forEach((function(e){r.removeAttribute(e)})))}))}},requires:["computeStyles"]};function h(e){return e.split("-")[0]}var m=Math.max,v=Math.min,g=Math.round;function y(){var e=navigator.userAgentData;return null!=e&&e.brands?e.brands.map((function(e){return e.brand+"/"+e.version})).join(" "):navigator.userAgent}function b(){return!/^((?!chrome|android).)*safari/i.test(y())}function w(e,t,n){void 0===t&&(t=!1),void 0===n&&(n=!1);var r=e.getBoundingClientRect(),o=1,i=1;t&&u(e)&&(o=e.offsetWidth>0&&g(r.width)/e.offsetWidth||1,i=e.offsetHeight>0&&g(r.height)/e.offsetHeight||1);var a=(p(e)?c(e):window).visualViewport,s=!b()&&n,f=(r.left+(s&&a?a.offsetLeft:0))/o,l=(r.top+(s&&a?a.offsetTop:0))/i,d=r.width/o,h=r.height/i;return{width:d,height:h,top:l,right:f+d,bottom:l+h,left:f,x:f,y:l}}function x(e){var t=w(e),n=e.offsetWidth,r=e.offsetHeight;return Math.abs(t.width-n)<=1&&(n=t.width),Math.abs(t.height-r)<=1&&(r=t.height),{x:e.offsetLeft,y:e.offsetTop,width:n,height:r}}function O(e,t){var n=t.getRootNode&&t.getRootNode();if(e.contains(t))return!0;if(n&&l(n)){var r=t;do{if(r&&e.isSameNode(r))return!0;r=r.parentNode||r.host}while(r)}return!1}function j(e){return c(e).getComputedStyle(e)}function E(e){return["table","td","th"].indexOf(f(e))>=0}function D(e){return((p(e)?e.ownerDocument:e.document)||window.document).documentElement}function A(e){return"html"===f(e)?e:e.assignedSlot||e.parentNode||(l(e)?e.host:null)||D(e)}function k(e){return u(e)&&"fixed"!==j(e).position?e.offsetParent:null}function L(e){for(var t=c(e),n=k(e);n&&E(n)&&"static"===j(n).position;)n=k(n);return n&&("html"===f(n)||"body"===f(n)&&"static"===j(n).position)?t:n||function(e){var t=/firefox/i.test(y());if(/Trident/i.test(y())&&u(e)&&"fixed"===j(e).position)return null;var n=A(e);for(l(n)&&(n=n.host);u(n)&&["html","body"].indexOf(f(n))<0;){var r=j(n);if("none"!==r.transform||"none"!==r.perspective||"paint"===r.contain||-1!==["transform","perspective"].indexOf(r.willChange)||t&&"filter"===r.willChange||t&&r.filter&&"none"!==r.filter)return n;n=n.parentNode}return null}(e)||t}function W(e){return["top","bottom"].indexOf(e)>=0?"x":"y"}function M(e,t,n){return m(e,v(t,n))}function P(e){return Object.assign({},{top:0,right:0,bottom:0,left:0},e)}function B(e,t){return t.reduce((function(t,n){return t[n]=e,t}),{})}function H(e){return e.split("-")[1]}var R={top:"auto",right:"auto",bottom:"auto",left:"auto"};function T(o){var i,a=o.popper,s=o.popperRect,f=o.placement,p=o.variation,u=o.offsets,l=o.position,d=o.gpuAcceleration,h=o.adaptive,m=o.roundOffsets,v=o.isFixed,y=u.x,b=void 0===y?0:y,w=u.y,x=void 0===w?0:w,O="function"==typeof m?m({x:b,y:x}):{x:b,y:x};b=O.x,x=O.y;var E=u.hasOwnProperty("x"),A=u.hasOwnProperty("y"),k=r,W=e,M=window;if(h){var P=L(a),B="clientHeight",H="clientWidth";if(P===c(a)&&"static"!==j(P=D(a)).position&&"absolute"===l&&(B="scrollHeight",H="scrollWidth"),f===e||(f===r||f===n)&&"end"===p)W=t,x-=(v&&P===M&&M.visualViewport?M.visualViewport.height:P[B])-s.height,x*=d?1:-1;if(f===r||(f===e||f===t)&&"end"===p)k=n,b-=(v&&P===M&&M.visualViewport?M.visualViewport.width:P[H])-s.width,b*=d?1:-1}var T,S=Object.assign({position:l},h&&R),V=!0===m?function(e){var t=e.x,n=e.y,r=window.devicePixelRatio||1;return{x:g(t*r)/r||0,y:g(n*r)/r||0}}({x:b,y:x}):{x:b,y:x};return b=V.x,x=V.y,d?Object.assign({},S,((T={})[W]=A?"0":"",T[k]=E?"0":"",T.transform=(M.devicePixelRatio||1)<=1?"translate("+b+"px, "+x+"px)":"translate3d("+b+"px, "+x+"px, 0)",T)):Object.assign({},S,((i={})[W]=A?x+"px":"",i[k]=E?b+"px":"",i.transform="",i))}var S={passive:!0};var V={left:"right",right:"left",bottom:"top",top:"bottom"};function q(e){return e.replace(/left|right|bottom|top/g,(function(e){return V[e]}))}var C={start:"end",end:"start"};function N(e){return e.replace(/start|end/g,(function(e){return C[e]}))}function I(e){var t=c(e);return{scrollLeft:t.pageXOffset,scrollTop:t.pageYOffset}}function F(e){return w(D(e)).left+I(e).scrollLeft}function U(e){var t=j(e),n=t.overflow,r=t.overflowX,o=t.overflowY;return/auto|scroll|overlay|hidden/.test(n+o+r)}function z(e){return["html","body","#document"].indexOf(f(e))>=0?e.ownerDocument.body:u(e)&&U(e)?e:z(A(e))}function _(e,t){var n;void 0===t&&(t=[]);var r=z(e),o=r===(null==(n=e.ownerDocument)?void 0:n.body),i=c(r),a=o?[i].concat(i.visualViewport||[],U(r)?r:[]):r,s=t.concat(a);return o?s:s.concat(_(A(a)))}function X(e){return Object.assign({},e,{left:e.x,top:e.y,right:e.x+e.width,bottom:e.y+e.height})}function Y(e,t,n){return"viewport"===t?X(function(e,t){var n=c(e),r=D(e),o=n.visualViewport,i=r.clientWidth,a=r.clientHeight,s=0,f=0;if(o){i=o.width,a=o.height;var p=b();(p||!p&&"fixed"===t)&&(s=o.offsetLeft,f=o.offsetTop)}return{width:i,height:a,x:s+F(e),y:f}}(e,n)):p(t)?function(e,t){var n=w(e,!1,"fixed"===t);return n.top=n.top+e.clientTop,n.left=n.left+e.clientLeft,n.bottom=n.top+e.clientHeight,n.right=n.left+e.clientWidth,n.width=e.clientWidth,n.height=e.clientHeight,n.x=n.left,n.y=n.top,n}(t,n):X(function(e){var t,n=D(e),r=I(e),o=null==(t=e.ownerDocument)?void 0:t.body,i=m(n.scrollWidth,n.clientWidth,o?o.scrollWidth:0,o?o.clientWidth:0),a=m(n.scrollHeight,n.clientHeight,o?o.scrollHeight:0,o?o.clientHeight:0),s=-r.scrollLeft+F(e),f=-r.scrollTop;return"rtl"===j(o||n).direction&&(s+=m(n.clientWidth,o?o.clientWidth:0)-i),{width:i,height:a,x:s,y:f}}(D(e)))}function G(e,t,n,r){var o="clippingParents"===t?function(e){var t=_(A(e)),n=["absolute","fixed"].indexOf(j(e).position)>=0&&u(e)?L(e):e;return p(n)?t.filter((function(e){return p(e)&&O(e,n)&&"body"!==f(e)})):[]}(e):[].concat(t),i=[].concat(o,[n]),a=i[0],s=i.reduce((function(t,n){var o=Y(e,n,r);return t.top=m(o.top,t.top),t.right=v(o.right,t.right),t.bottom=v(o.bottom,t.bottom),t.left=m(o.left,t.left),t}),Y(e,a,r));return s.width=s.right-s.left,s.height=s.bottom-s.top,s.x=s.left,s.y=s.top,s}function J(o){var i,a=o.reference,s=o.element,f=o.placement,c=f?h(f):null,p=f?H(f):null,u=a.x+a.width/2-s.width/2,l=a.y+a.height/2-s.height/2;switch(c){case e:i={x:u,y:a.y-s.height};break;case t:i={x:u,y:a.y+a.height};break;case n:i={x:a.x+a.width,y:l};break;case r:i={x:a.x-s.width,y:l};break;default:i={x:a.x,y:a.y}}var d=c?W(c):null;if(null!=d){var m="y"===d?"height":"width";switch(p){case"start":i[d]=i[d]-(a[m]/2-s[m]/2);break;case"end":i[d]=i[d]+(a[m]/2-s[m]/2)}}return i}function K(r,i){void 0===i&&(i={});var a=i,s=a.placement,f=void 0===s?r.placement:s,c=a.strategy,u=void 0===c?r.strategy:c,l=a.boundary,d=void 0===l?"clippingParents":l,h=a.rootBoundary,m=void 0===h?"viewport":h,v=a.elementContext,g=void 0===v?"popper":v,y=a.altBoundary,b=void 0!==y&&y,x=a.padding,O=void 0===x?0:x,j=P("number"!=typeof O?O:B(O,o)),E="popper"===g?"reference":"popper",A=r.rects.popper,k=r.elements[b?E:g],L=G(p(k)?k:k.contextElement||D(r.elements.popper),d,m,u),W=w(r.elements.reference),M=J({reference:W,element:A,strategy:"absolute",placement:f}),H=X(Object.assign({},A,M)),R="popper"===g?H:W,T={top:L.top-R.top+j.top,bottom:R.bottom-L.bottom+j.bottom,left:L.left-R.left+j.left,right:R.right-L.right+j.right},S=r.modifiersData.offset;if("popper"===g&&S){var V=S[f];Object.keys(T).forEach((function(r){var o=[n,t].indexOf(r)>=0?1:-1,i=[e,t].indexOf(r)>=0?"y":"x";T[r]+=V[i]*o}))}return T}function Q(e,t,n){return void 0===n&&(n={x:0,y:0}),{top:e.top-t.height-n.y,right:e.right-t.width+n.x,bottom:e.bottom-t.height+n.y,left:e.left-t.width-n.x}}function Z(o){return[e,n,t,r].some((function(e){return o[e]>=0}))}function $(e,t,n){void 0===n&&(n=!1);var r,o,i=u(t),a=u(t)&&function(e){var t=e.getBoundingClientRect(),n=g(t.width)/e.offsetWidth||1,r=g(t.height)/e.offsetHeight||1;return 1!==n||1!==r}(t),s=D(t),p=w(e,a,n),l={scrollLeft:0,scrollTop:0},d={x:0,y:0};return(i||!i&&!n)&&(("body"!==f(t)||U(s))&&(l=(r=t)!==c(r)&&u(r)?{scrollLeft:(o=r).scrollLeft,scrollTop:o.scrollTop}:I(r)),u(t)?((d=w(t,!0)).x+=t.clientLeft,d.y+=t.clientTop):s&&(d.x=F(s))),{x:p.left+l.scrollLeft-d.x,y:p.top+l.scrollTop-d.y,width:p.width,height:p.height}}function ee(e){var t=new Map,n=new Set,r=[];function o(e){n.add(e.name),[].concat(e.requires||[],e.requiresIfExists||[]).forEach((function(e){if(!n.has(e)){var r=t.get(e);r&&o(r)}})),r.push(e)}return e.forEach((function(e){t.set(e.name,e)})),e.forEach((function(e){n.has(e.name)||o(e)})),r}var te={placement:"bottom",modifiers:[],strategy:"absolute"};function ne(){for(var e=arguments.length,t=new Array(e),n=0;n<e;n++)t[n]=arguments[n];return!t.some((function(e){return!(e&&"function"==typeof e.getBoundingClientRect)}))}function re(e){void 0===e&&(e={});var t=e,n=t.defaultModifiers,r=void 0===n?[]:n,o=t.defaultOptions,i=void 0===o?te:o;return function(e,t,n){void 0===n&&(n=i);var o,a,f={placement:"bottom",orderedModifiers:[],options:Object.assign({},te,i),modifiersData:{},elements:{reference:e,popper:t},attributes:{},styles:{}},c=[],u=!1,l={state:f,setOptions:function(n){var o="function"==typeof n?n(f.options):n;d(),f.options=Object.assign({},i,f.options,o),f.scrollParents={reference:p(e)?_(e):e.contextElement?_(e.contextElement):[],popper:_(t)};var a,u,h=function(e){var t=ee(e);return s.reduce((function(e,n){return e.concat(t.filter((function(e){return e.phase===n})))}),[])}((a=[].concat(r,f.options.modifiers),u=a.reduce((function(e,t){var n=e[t.name];return e[t.name]=n?Object.assign({},n,t,{options:Object.assign({},n.options,t.options),data:Object.assign({},n.data,t.data)}):t,e}),{}),Object.keys(u).map((function(e){return u[e]}))));return f.orderedModifiers=h.filter((function(e){return e.enabled})),f.orderedModifiers.forEach((function(e){var t=e.name,n=e.options,r=void 0===n?{}:n,o=e.effect;if("function"==typeof o){var i=o({state:f,name:t,instance:l,options:r}),a=function(){};c.push(i||a)}})),l.update()},forceUpdate:function(){if(!u){var e=f.elements,t=e.reference,n=e.popper;if(ne(t,n)){f.rects={reference:$(t,L(n),"fixed"===f.options.strategy),popper:x(n)},f.reset=!1,f.placement=f.options.placement,f.orderedModifiers.forEach((function(e){return f.modifiersData[e.name]=Object.assign({},e.data)}));for(var r=0;r<f.orderedModifiers.length;r++)if(!0!==f.reset){var o=f.orderedModifiers[r],i=o.fn,a=o.options,s=void 0===a?{}:a,c=o.name;"function"==typeof i&&(f=i({state:f,options:s,name:c,instance:l})||f)}else f.reset=!1,r=-1}}},update:(o=function(){return new Promise((function(e){l.forceUpdate(),e(f)}))},function(){return a||(a=new Promise((function(e){Promise.resolve().then((function(){a=void 0,e(o())}))}))),a}),destroy:function(){d(),u=!0}};if(!ne(e,t))return l;function d(){c.forEach((function(e){return e()})),c=[]}return l.setOptions(n).then((function(e){!u&&n.onFirstUpdate&&n.onFirstUpdate(e)})),l}}var oe=re({defaultModifiers:[{name:"eventListeners",enabled:!0,phase:"write",fn:function(){},effect:function(e){var t=e.state,n=e.instance,r=e.options,o=r.scroll,i=void 0===o||o,a=r.resize,s=void 0===a||a,f=c(t.elements.popper),p=[].concat(t.scrollParents.reference,t.scrollParents.popper);return i&&p.forEach((function(e){e.addEventListener("scroll",n.update,S)})),s&&f.addEventListener("resize",n.update,S),function(){i&&p.forEach((function(e){e.removeEventListener("scroll",n.update,S)})),s&&f.removeEventListener("resize",n.update,S)}},data:{}},{name:"popperOffsets",enabled:!0,phase:"read",fn:function(e){var t=e.state,n=e.name;t.modifiersData[n]=J({reference:t.rects.reference,element:t.rects.popper,strategy:"absolute",placement:t.placement})},data:{}},{name:"computeStyles",enabled:!0,phase:"beforeWrite",fn:function(e){var t=e.state,n=e.options,r=n.gpuAcceleration,o=void 0===r||r,i=n.adaptive,a=void 0===i||i,s=n.roundOffsets,f=void 0===s||s,c={placement:h(t.placement),variation:H(t.placement),popper:t.elements.popper,popperRect:t.rects.popper,gpuAcceleration:o,isFixed:"fixed"===t.options.strategy};null!=t.modifiersData.popperOffsets&&(t.styles.popper=Object.assign({},t.styles.popper,T(Object.assign({},c,{offsets:t.modifiersData.popperOffsets,position:t.options.strategy,adaptive:a,roundOffsets:f})))),null!=t.modifiersData.arrow&&(t.styles.arrow=Object.assign({},t.styles.arrow,T(Object.assign({},c,{offsets:t.modifiersData.arrow,position:"absolute",adaptive:!1,roundOffsets:f})))),t.attributes.popper=Object.assign({},t.attributes.popper,{"data-popper-placement":t.placement})},data:{}},d,{name:"offset",enabled:!0,phase:"main",requires:["popperOffsets"],fn:function(t){var o=t.state,i=t.options,s=t.name,f=i.offset,c=void 0===f?[0,0]:f,p=a.reduce((function(t,i){return t[i]=function(t,o,i){var a=h(t),s=[r,e].indexOf(a)>=0?-1:1,f="function"==typeof i?i(Object.assign({},o,{placement:t})):i,c=f[0],p=f[1];return c=c||0,p=(p||0)*s,[r,n].indexOf(a)>=0?{x:p,y:c}:{x:c,y:p}}(i,o.rects,c),t}),{}),u=p[o.placement],l=u.x,d=u.y;null!=o.modifiersData.popperOffsets&&(o.modifiersData.popperOffsets.x+=l,o.modifiersData.popperOffsets.y+=d),o.modifiersData[s]=p}},{name:"flip",enabled:!0,phase:"main",fn:function(s){var f=s.state,c=s.options,p=s.name;if(!f.modifiersData[p]._skip){for(var u=c.mainAxis,l=void 0===u||u,d=c.altAxis,m=void 0===d||d,v=c.fallbackPlacements,g=c.padding,y=c.boundary,b=c.rootBoundary,w=c.altBoundary,x=c.flipVariations,O=void 0===x||x,j=c.allowedAutoPlacements,E=f.options.placement,D=h(E),A=v||(D===E||!O?[q(E)]:function(e){if("auto"===h(e))return[];var t=q(e);return[N(e),t,N(t)]}(E)),k=[E].concat(A).reduce((function(e,t){return e.concat("auto"===h(t)?function(e,t){void 0===t&&(t={});var n=t,r=n.placement,s=n.boundary,f=n.rootBoundary,c=n.padding,p=n.flipVariations,u=n.allowedAutoPlacements,l=void 0===u?a:u,d=H(r),m=d?p?i:i.filter((function(e){return H(e)===d})):o,v=m.filter((function(e){return l.indexOf(e)>=0}));0===v.length&&(v=m);var g=v.reduce((function(t,n){return t[n]=K(e,{placement:n,boundary:s,rootBoundary:f,padding:c})[h(n)],t}),{});return Object.keys(g).sort((function(e,t){return g[e]-g[t]}))}(f,{placement:t,boundary:y,rootBoundary:b,padding:g,flipVariations:O,allowedAutoPlacements:j}):t)}),[]),L=f.rects.reference,W=f.rects.popper,M=new Map,P=!0,B=k[0],R=0;R<k.length;R++){var T=k[R],S=h(T),V="start"===H(T),C=[e,t].indexOf(S)>=0,I=C?"width":"height",F=K(f,{placement:T,boundary:y,rootBoundary:b,altBoundary:w,padding:g}),U=C?V?n:r:V?t:e;L[I]>W[I]&&(U=q(U));var z=q(U),_=[];if(l&&_.push(F[S]<=0),m&&_.push(F[U]<=0,F[z]<=0),_.every((function(e){return e}))){B=T,P=!1;break}M.set(T,_)}if(P)for(var X=function(e){var t=k.find((function(t){var n=M.get(t);if(n)return n.slice(0,e).every((function(e){return e}))}));if(t)return B=t,"break"},Y=O?3:1;Y>0;Y--){if("break"===X(Y))break}f.placement!==B&&(f.modifiersData[p]._skip=!0,f.placement=B,f.reset=!0)}},requiresIfExists:["offset"],data:{_skip:!1}},{name:"preventOverflow",enabled:!0,phase:"main",fn:function(o){var i=o.state,a=o.options,s=o.name,f=a.mainAxis,c=void 0===f||f,p=a.altAxis,u=void 0!==p&&p,l=a.boundary,d=a.rootBoundary,g=a.altBoundary,y=a.padding,b=a.tether,w=void 0===b||b,O=a.tetherOffset,j=void 0===O?0:O,E=K(i,{boundary:l,rootBoundary:d,padding:y,altBoundary:g}),D=h(i.placement),A=H(i.placement),k=!A,P=W(D),B="x"===P?"y":"x",R=i.modifiersData.popperOffsets,T=i.rects.reference,S=i.rects.popper,V="function"==typeof j?j(Object.assign({},i.rects,{placement:i.placement})):j,q="number"==typeof V?{mainAxis:V,altAxis:V}:Object.assign({mainAxis:0,altAxis:0},V),C=i.modifiersData.offset?i.modifiersData.offset[i.placement]:null,N={x:0,y:0};if(R){if(c){var I,F="y"===P?e:r,U="y"===P?t:n,z="y"===P?"height":"width",_=R[P],X=_+E[F],Y=_-E[U],G=w?-S[z]/2:0,J="start"===A?T[z]:S[z],Q="start"===A?-S[z]:-T[z],Z=i.elements.arrow,$=w&&Z?x(Z):{width:0,height:0},ee=i.modifiersData["arrow#persistent"]?i.modifiersData["arrow#persistent"].padding:{top:0,right:0,bottom:0,left:0},te=ee[F],ne=ee[U],re=M(0,T[z],$[z]),oe=k?T[z]/2-G-re-te-q.mainAxis:J-re-te-q.mainAxis,ie=k?-T[z]/2+G+re+ne+q.mainAxis:Q+re+ne+q.mainAxis,ae=i.elements.arrow&&L(i.elements.arrow),se=ae?"y"===P?ae.clientTop||0:ae.clientLeft||0:0,fe=null!=(I=null==C?void 0:C[P])?I:0,ce=_+ie-fe,pe=M(w?v(X,_+oe-fe-se):X,_,w?m(Y,ce):Y);R[P]=pe,N[P]=pe-_}if(u){var ue,le="x"===P?e:r,de="x"===P?t:n,he=R[B],me="y"===B?"height":"width",ve=he+E[le],ge=he-E[de],ye=-1!==[e,r].indexOf(D),be=null!=(ue=null==C?void 0:C[B])?ue:0,we=ye?ve:he-T[me]-S[me]-be+q.altAxis,xe=ye?he+T[me]+S[me]-be-q.altAxis:ge,Oe=w&&ye?(Ee=M(we,he,je=xe))>je?je:Ee:M(w?we:ve,he,w?xe:ge);R[B]=Oe,N[B]=Oe-he}var je,Ee;i.modifiersData[s]=N}},requiresIfExists:["offset"]},{name:"arrow",enabled:!0,phase:"main",fn:function(i){var a,s=i.state,f=i.name,c=i.options,p=s.elements.arrow,u=s.modifiersData.popperOffsets,l=h(s.placement),d=W(l),m=[r,n].indexOf(l)>=0?"height":"width";if(p&&u){var v=function(e,t){return P("number"!=typeof(e="function"==typeof e?e(Object.assign({},t.rects,{placement:t.placement})):e)?e:B(e,o))}(c.padding,s),g=x(p),y="y"===d?e:r,b="y"===d?t:n,w=s.rects.reference[m]+s.rects.reference[d]-u[d]-s.rects.popper[m],O=u[d]-s.rects.reference[d],j=L(p),E=j?"y"===d?j.clientHeight||0:j.clientWidth||0:0,D=w/2-O/2,A=v[y],k=E-g[m]-v[b],H=E/2-g[m]/2+D,R=M(A,H,k),T=d;s.modifiersData[f]=((a={})[T]=R,a.centerOffset=R-H,a)}},effect:function(e){var t=e.state,n=e.options.element,r=void 0===n?"[data-popper-arrow]":n;null!=r&&("string"!=typeof r||(r=t.elements.popper.querySelector(r)))&&O(t.elements.popper,r)&&(t.elements.arrow=r)},requires:["popperOffsets"],requiresIfExists:["preventOverflow"]},{name:"hide",enabled:!0,phase:"main",requiresIfExists:["preventOverflow"],fn:function(e){var t=e.state,n=e.name,r=t.rects.reference,o=t.rects.popper,i=t.modifiersData.preventOverflow,a=K(t,{elementContext:"reference"}),s=K(t,{altBoundary:!0}),f=Q(a,r),c=Q(s,o,i),p=Z(f),u=Z(c);t.modifiersData[n]={referenceClippingOffsets:f,popperEscapeOffsets:c,isReferenceHidden:p,hasPopperEscaped:u},t.attributes.popper=Object.assign({},t.attributes.popper,{"data-popper-reference-hidden":p,"data-popper-escaped":u})}}]});export{oe as c};
