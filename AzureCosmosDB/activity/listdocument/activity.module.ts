import { HttpModule } from "@angular/http";
import { NgModule } from "@angular/core";
import { CommonModule } from "@angular/common";

import { ListDocumentActivityContributionHandler} from "./activity";
import { WiServiceContribution} from "wi-studio/app/contrib/wi-contrib";

@NgModule({
  imports: [
    CommonModule,
    HttpModule
  ],
  exports: [],
  declarations: [],
  entryComponents: [],
  providers: [
    {
       provide: WiServiceContribution,
       useClass: ListDocumentActivityContributionHandler
     }
  ],
  bootstrap: []
})

export default class ListDocumentActivityModule {

}