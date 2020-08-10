"use strict";var __extends=this&&this.__extends||function(){var e=function(t,r){return(e=Object.setPrototypeOf||{__proto__:[]}instanceof Array&&function(e,t){e.__proto__=t}||function(e,t){for(var r in t)t.hasOwnProperty(r)&&(e[r]=t[r])})(t,r)};return function(t,r){function o(){this.constructor=t}e(t,r),t.prototype=null===r?Object.create(r):(o.prototype=r.prototype,new o)}}(),__decorate=this&&this.__decorate||function(e,t,r,o){var n,i=arguments.length,l=i<3?t:null===o?o=Object.getOwnPropertyDescriptor(t,r):o;if("object"==typeof Reflect&&"function"==typeof Reflect.decorate)l=Reflect.decorate(e,t,r,o);else for(var a=e.length-1;a>=0;a--)(n=e[a])&&(l=(i<3?n(l):i>3?n(t,r,l):n(t,r))||l);return i>3&&l&&Object.defineProperty(t,r,l),l},__metadata=this&&this.__metadata||function(e,t){if("object"==typeof Reflect&&"function"==typeof Reflect.metadata)return Reflect.metadata(e,t)};Object.defineProperty(exports,"__esModule",{value:!0});var core_1=require("@angular/core"),http_1=require("@angular/http"),Observable_1=require("rxjs/Observable"),wi_contrib_1=require("wi-studio/app/contrib/wi-contrib"),lodash=require("lodash"),FUnitMockTriggerHandler=function(e){function t(t,r,o){var n=e.call(this,t,r,o)||this;return n.injector=t,n.http=r,n.contribModelService=o,n.value=function(e,t){var r=n.getModelService().getApplication();switch(e){case"flowURI":var o=[],i=[];if(r){r.getTriggerFlowModelMaps().map(function(e){console.log("trigger name = "+e.getTriggerElement().name),console.log("trigger flow = "+e.getFlowModel().getName()),console.log("isTriggerFlow = "+e.getFlowModel().isTriggerFlow()),console.log("current flow= "+t.getCurrentFlowName()),"funit"==e.getTriggerElement().name||"mock"==e.getTriggerElement().name||t.getCurrentFlowName()===e.getFlowModel().getName()||e.getFlowModel().isTriggerFlow()||-1===o.indexOf(e.getFlowModel().getName())&&(o.push(e.getFlowModel().getName()),i.push({unique_id:"res://flow:"+e.getFlowModel().getName().replace(/ /g,"_"),name:e.getFlowModel().getName()}))})}return i}return null},n.validate=function(e,t){return null},n.action=function(e,t){var r=wi_contrib_1.CreateFlowActionResult.newActionResult();return Observable_1.Observable.create(function(e){n.createFlow(t,r);var o=wi_contrib_1.ActionResult.newActionResult().setSuccess(!0).setResult(r);e.next(o)})},n}return __extends(t,e),t.prototype.createFlow=function(e,t){var r=this.getModelService(),o=(r.getApplication(),e.getField("flowURI").value),n=r.createTriggerElement("funit/mock");if(n)for(var i=0;i<n.handler.settings.length;i++)"flowURI"===n.handler.settings[i].name&&(n.handler.settings[i].value=o);var l=e.getFlowName(),a=r.createFlow(l,e.getFlowDescription());t=t.addTriggerFlowMapping(lodash.cloneDeep(n),lodash.cloneDeep(a))},t=__decorate([wi_contrib_1.WiContrib({}),core_1.Injectable(),__metadata("design:paramtypes",[core_1.Injector,http_1.Http,wi_contrib_1.WiContribModelService])],t)}(wi_contrib_1.WiServiceHandlerContribution);exports.FUnitMockTriggerHandler=FUnitMockTriggerHandler;