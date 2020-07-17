import { Injectable, Injector } from "@angular/core";
import { WiContrib, WiServiceHandlerContribution, WiContributionUtils,AUTHENTICATION_TYPE } from "wi-studio/app/contrib/wi-contrib";
import { IConnectorContribution } from "wi-studio/common/models/contrib";
import { IValidationResult, ValidationResult, ValidationError } from "wi-studio/common/models/validation";
import { IActionResult, ActionResult } from "wi-studio/common/models/contrib";
import { Observable } from "rxjs/Observable";
import {Http} from "@angular/http";

@Injectable()
@WiContrib({})
export class ConnectorConnectorService extends WiServiceHandlerContribution {
    constructor( private injector: Injector, private http: Http) {
        super(injector, http);
    }

    value = (fieldName: string, context: IConnectorContribution): Observable<any> | any => {
       switch(fieldName){
           case "connectorType":
               let cnn = []
               cnn.push({
                   "unique_id": "json-api",
                    "name": "Digital Asset Json Api Server"})
                return cnn
       }
        return null;    
    }

    validate = (fieldName: string, context: IConnectorContribution): Observable<IValidationResult> | IValidationResult => {
       let cnn = context.getField("connectorType").value
       if(Boolean(cnn)){
           switch(fieldName){
                case "host":
                case "port":
                case "JWT":
                case "timeout":
                    return Observable.create(observer => {
                        let vresult: IValidationResult = ValidationResult.newValidationResult();
                        vresult.setVisible(true);
                        observer.next(vresult);
                    });
           }
       }
        return null;
    }

    action = (name: string, context: IConnectorContribution): Observable<IActionResult> | IActionResult => {

       if(name === "Done"){
      
            return Observable.create(observer => {
                let actionResult = {
                    context: context,
                    authType: AUTHENTICATION_TYPE.BASIC,
                    authData: {}
                }
                observer.next(ActionResult.newActionResult().setResult(actionResult));
            })
        } 
        return null
          
    }
}
