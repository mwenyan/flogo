"use strict";var __decorate=this&&this.__decorate||function(t,e,o,r){var i,c=arguments.length,n=c<3?e:null===r?r=Object.getOwnPropertyDescriptor(e,o):r;if("object"==typeof Reflect&&"function"==typeof Reflect.decorate)n=Reflect.decorate(t,e,o,r);else for(var u=t.length-1;u>=0;u--)(i=t[u])&&(n=(c<3?i(n):c>3?i(e,o,n):i(e,o))||n);return c>3&&n&&Object.defineProperty(e,o,n),n};Object.defineProperty(exports,"__esModule",{value:!0});var http_1=require("@angular/http"),core_1=require("@angular/core"),common_1=require("@angular/common"),activity_1=require("./activity"),wi_contrib_1=require("wi-studio/app/contrib/wi-contrib"),ListDocumentActivityModule=function(){function t(){}return t=__decorate([core_1.NgModule({imports:[common_1.CommonModule,http_1.HttpModule],exports:[],declarations:[],entryComponents:[],providers:[{provide:wi_contrib_1.WiServiceContribution,useClass:activity_1.ListDocumentActivityContributionHandler}],bootstrap:[]})],t)}();exports.default=ListDocumentActivityModule;