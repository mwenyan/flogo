import { HttpModule, Http } from "@angular/http";
import { NgModule } from "@angular/core";
import { SendMailActivityService, UserNameFieldProvider, PasswordFieldProvider } from "./activity";
import { WiServiceContribution } from "wi-studio/app/contrib/wi-contrib";


@NgModule({
  imports: [
    HttpModule,
  ],
  exports: [

  ],
  declarations: [

  ],
  entryComponents: [

  ],
  providers: [
    {
       provide: WiServiceContribution,
       useClass: SendMailActivityService
     },
     {
       provide: UserNameFieldProvider,
       useClass: UserNameFieldProvider
     },
     {
       provide: PasswordFieldProvider,
       useClass: PasswordFieldProvider
     }
  ],
  bootstrap: []
})
export default class SendMailContribModule {

}
