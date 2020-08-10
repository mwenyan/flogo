import {NgModule} from "@angular/core"
import {HttpModule} from "@angular/http";
import {WiServiceContribution} from "wi-studio/app/contrib/wi-contrib"
import {FUnitMockTriggerHandler} from "./trigger"

@NgModule({
    imports: [
        HttpModule
    ],
    providers: [
        {
            provide: WiServiceContribution,
            useClass: FUnitMockTriggerHandler
        }
    ]
})
export default class FUnitMockTriggerModule {

}
