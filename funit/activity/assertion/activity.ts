import {Observable} from "rxjs/Observable";
import {Injectable, Injector, Inject} from "@angular/core";
import {Http} from "@angular/http";
import {
    WiContrib,
    WiServiceHandlerContribution,
    IValidationResult,
    ValidationResult,
    IActivityContribution,
    WiContributionUtils,
    IConnectorContribution
} from "wi-studio/app/contrib/wi-contrib";

import * as lodash from "lodash";

@WiContrib({})
@Injectable()
export class AssertionActivityContributionHandler extends WiServiceHandlerContribution {
    constructor(private injector: Injector, private http: Http,) {
        super(injector, http);
    }
   
    value = (fieldName: string, context: IActivityContribution): any | Observable<any> => {

        switch(fieldName) {
            
            case "input":
                let msg = context.getField("msg").value
                let count = 0
                if (Boolean(msg)) {
                    count = (msg.match(/%v/g) || []).length;
                }

                let schema = {
                    "$schema": "http://json-schema.org/draft-04/schema#",
                    "type": "object",
                    "properties": {
                        "assertion": {
                            "type": "boolean"
                        }
                    },
                    "required":["assertion"]
                }

                for(var i = 0; i < count; i++){
                    schema.properties["msg_var" + i] = {"type": "any"}
                    schema.required.push("msg_var" + i)
                }
               
                 
                return JSON.stringify(schema);
            case "output":
                let oschemaobj = {assertResult: true, message: "" }
                  
                return JSON.stringify(oschemaobj)
                
        }
        return null;
    }
 
    validate = (fieldName: string, context: IActivityContribution): Observable<IValidationResult> | IValidationResult => {
        return null;
    }
}