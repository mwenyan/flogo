"use strict";var __extends=this&&this.__extends||function(){var e=Object.setPrototypeOf||{__proto__:[]}instanceof Array&&function(e,t){e.__proto__=t}||function(e,t){for(var n in t)t.hasOwnProperty(n)&&(e[n]=t[n])};return function(t,n){function r(){this.constructor=t}e(t,n),t.prototype=null===n?Object.create(n):(r.prototype=n.prototype,new r)}}(),__decorate=this&&this.__decorate||function(e,t,n,r){var o,i=arguments.length,s=i<3?t:null===r?r=Object.getOwnPropertyDescriptor(t,n):r;if("object"==typeof Reflect&&"function"==typeof Reflect.decorate)s=Reflect.decorate(e,t,n,r);else for(var a=e.length-1;a>=0;a--)(o=e[a])&&(s=(i<3?o(s):i>3?o(t,n,s):o(t,n))||s);return i>3&&s&&Object.defineProperty(t,n,s),s},__metadata=this&&this.__metadata||function(e,t){if("object"==typeof Reflect&&"function"==typeof Reflect.metadata)return Reflect.metadata(e,t)},__param=this&&this.__param||function(e,t){return function(n,r){t(n,r,e)}};Object.defineProperty(exports,"__esModule",{value:!0});var http_1=require("@angular/http"),core_1=require("@angular/core"),wi_contrib_1=require("wi-studio/app/contrib/wi-contrib"),contrib_1=require("wi-studio/common/models/contrib"),rxjs_extensions_1=require("wi-studio/common/rxjs-extensions");require("wi-studio/common/rxjs-extensions");var validation_1=require("wi-studio/common/models/validation");require("rxjs/add/observable/zip");var cryptos=require("crypto-js"),schemadoc_1=require("./schemadoc"),Result=function(e,t,n){},Connection=function(){function e(e){var t=this;this.http=e,this.oAuthProperties=function(){return rxjs_extensions_1.Observable.zip(rxjs_extensions_1.Observable.of(t.authorizationRuleName),rxjs_extensions_1.Observable.of(t.primarysecondaryKey),function(e,t){return{authorizationRuleName:e,primarysecondaryKey:t}})}}return e.getConnection=function(t,n){var r=new e(n);return rxjs_extensions_1.Observable.from(t).reduce(function(e,t){return e[t.name]=t.value,e},r)},e}();Connection.dateRegx=/^\d{4}-(?:0[1-9]|1[0-2])-(?:0[1-9]|[1-2]\d|3[0-1])T(?:[0-1]\d|2[0-3]):[0-5]\d:[0-5]\dZ/,exports.Connection=Connection;var TibcoAzServiceBusConnectorContribution=function(e){function t(t,n){var r=e.call(this,t,n)||this;return r.http=n,r.value=function(e,t){return r.validations[t.title]&&r.validations[t.title][e]&&delete r.validations[t.title][e],r.validations[t.title]||(r.validations[t.title]={}),null},r.validate=function(e,t){return null},r.connection=function(e){var t=new Connection(r.http);return rxjs_extensions_1.Observable.from(e).reduce(function(e,t){return e[t.name]=t.value,e},t)},r.isDuplicate=function(e,t){return wi_contrib_1.WiContributionUtils.getConnections(r.http,r.category).mergeMap(function(e){return e}).map(function(e){return{id:wi_contrib_1.WiContributionUtils.getUniqueId(e),ction:r.connection(e.settings)}}).filter(function(e){return e.id!==t}).mergeMap(function(e){return e.ction}).filter(function(t){return t.name===e.name}).reduce(function(e,t){return e.push(t),e},[]).map(function(e){return e.length>0})},r.defaultResult=function(e){var t={context:e,authType:wi_contrib_1.AUTHENTICATION_TYPE.BASIC,authData:{}};return rxjs_extensions_1.Observable.of(contrib_1.ActionResult.newActionResult().setResult(t))},r.action=function(e,t){var n;if("Login"===e)return r.connection(t.settings).switchMap(function(e){return r.isDuplicate(e,wi_contrib_1.WiContributionUtils.getUniqueId(t)).map(function(t){if(t)throw new Error("Connection with name "+e.name+" already exists");return e})}).switchMap(function(e){if(!(null!==e.name&&e.name.trim().length>0))throw new Error("Please enter the Connection Name!");if(""===e.resourceURI)throw new Error("Please enter the resource URI");var o="";if(o="https://"+e.resourceURI+".servicebus.windows.net",""===e.authorizationRuleName)throw new Error("Please enter the Authorization Rule Name!");if(""===e.primarysecondaryKey)throw new Error("Please enter the primary/secondaryKey!");var i="";i=null!==o&&o.trim().length>0?i+o:"errorMsg";var s,a,c=String((new Date).getTime()/1e3+86400);if(i=encodeURIComponent(i)+"\n"+c,console.log(i),i.includes("errorMsg"))throw new Error("errorMsg");if(!(null!==e.authorizationRuleName&&e.authorizationRuleName.trim().length>0))throw new Error("errorMsg");var u=e.primarysecondaryKey;return s=cryptos.HmacSHA256(i,u),a=cryptos.enc.Base64.stringify(s),a=encodeURIComponent(a),o=encodeURIComponent(o),console.log(o),e.primarysecondaryKey="",n="SharedAccessSignature sr="+o+"&sig="+a+"&se="+c+"&skn="+e.authorizationRuleName,o="https://"+e.resourceURI+".servicebus.windows.net",wi_contrib_1.WiProxyCORSUtils.createRequest(r.http,o+"/testconnectionqueue"+c).addHeader("Content-Type","application/atom+xml;type=entry;charset=utf-8").addHeader("Authorization",n).addHeader("If-Match","*").addMethod("PUT").addBody('<entry xmlns="http://www.w3.org/2005/Atom"><content type="application/xml"><QueueDescription xmlns="http://schemas.microsoft.com/netservices/2010/10/servicebus/connect"><MaxDeliveryCount>10</MaxDeliveryCount></QueueDescription></content></entry>').send().switchMap(function(e){if(e.status>=500||401===e.status)return rxjs_extensions_1.Observable.of(contrib_1.ActionResult.newActionResult().setSuccess(!1).setResult(new validation_1.ValidationError("AZSERVICEBUSCONN-1002","Connection Authentication error: "+e.statusText+": Check your connection parameters")));for(var n=0;n<t.settings.length;n++)"DocsMetadata"===t.settings[n].name&&(t.settings[n].value=schemadoc_1.JsonSchema.Types.schemaDoc());var r={context:t,authType:wi_contrib_1.AUTHENTICATION_TYPE.BASIC,authData:{}};return rxjs_extensions_1.Observable.of(contrib_1.ActionResult.newActionResult().setResult(r))})}).catch(function(e){if(e instanceof http_1.Response){if(e.status){if(e.status>=500||401===e.status)return rxjs_extensions_1.Observable.of(contrib_1.ActionResult.newActionResult().setSuccess(!1).setResult(new validation_1.ValidationError("AZSERVICEBUSCONN-1002","Connection Authentication error: "+e.statusText+": Check your connection parameters")));for(var n=0;n<t.settings.length;n++)"DocsMetadata"===t.settings[n].name&&(t.settings[n].value=schemadoc_1.JsonSchema.Types.schemaDoc());var r={context:t,authType:wi_contrib_1.AUTHENTICATION_TYPE.BASIC,authData:{}};return rxjs_extensions_1.Observable.of(contrib_1.ActionResult.newActionResult().setResult(r))}return rxjs_extensions_1.Observable.of(contrib_1.ActionResult.newActionResult().setSuccess(!1).setResult(new validation_1.ValidationError("AZSERVICEBUSCONN-1002","Connection Authentication error: "+e.statusText+": Check your connection parameters")))}if(e instanceof Error)return console.log("AuthenticationFailed: "+e),rxjs_extensions_1.Observable.of(contrib_1.ActionResult.newActionResult().setSuccess(!1).setResult(new validation_1.ValidationError("AZSERVICEBUSCONN-1005","Connection Authentication error: "+e.message)))})},r.validations={},r.category="AzureServiceBus",r}return __extends(t,e),t}(wi_contrib_1.WiServiceHandlerContribution);TibcoAzServiceBusConnectorContribution=__decorate([wi_contrib_1.WiContrib({}),core_1.Injectable(),__param(0,core_1.Inject(core_1.Injector)),__metadata("design:paramtypes",[Object,http_1.Http])],TibcoAzServiceBusConnectorContribution),exports.TibcoAzServiceBusConnectorContribution=TibcoAzServiceBusConnectorContribution;