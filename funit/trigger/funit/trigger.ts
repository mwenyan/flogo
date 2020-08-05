/// <amd-dependency path="./common"/>
import {Injectable, Injector} from "@angular/core";
import {Http} from "@angular/http";
import {Observable} from "rxjs/Observable";
import {
    WiContrib,
    WiServiceHandlerContribution,
    IValidationResult,
    ValidationResult,
    ActionResult,
    IActionResult,
    ICreateFlowActionContext,
    CreateFlowActionResult,
    WiContribModelService,
    WiContributionUtils
} from "wi-studio/app/contrib/wi-contrib";
import { ITriggerContribution, IFieldDefinition, IConnectionAllowedValue, MODE } from "wi-studio/common/models/contrib";

@WiContrib({})
@Injectable()
export class FUnitTriggerHandler extends WiServiceHandlerContribution {
    constructor(private injector: Injector, private http: Http, private contribModelService: WiContribModelService) {
        super(injector, http, contribModelService);
    }
    
    value = (fieldName: string, context: ITriggerContribution): Observable<any> | any => {
      /*  switch(fieldName) {
            case "report":
                let oschemaobj = [{assertion: "", pass: true, message: "" }]
                  
                return JSON.stringify(oschemaobj)
        }*/
        return null
    }

    validate = (fieldName: string, context: ITriggerContribution): Observable<IValidationResult> | IValidationResult => {
        return null
    }

    action = (actionId: string, context: ICreateFlowActionContext): Observable<IActionResult> | IActionResult => {
        let result = CreateFlowActionResult.newActionResult();
        return Observable.create(observer => {                                    
            let actionResult = ActionResult.newActionResult().setSuccess(true).setResult(result);
            observer.next(actionResult);
        });
    }
}