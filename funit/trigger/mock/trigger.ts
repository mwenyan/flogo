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
import * as lodash from "lodash";

@WiContrib({})
@Injectable()
export class FUnitMockTriggerHandler extends WiServiceHandlerContribution {
    constructor(private injector: Injector, private http: Http, private contribModelService: WiContribModelService) {
        super(injector, http, contribModelService);
    }
    
    value = (fieldName: string, context: ITriggerContribution): Observable<any> | any => {
        var modelService = this.getModelService();
        var applicationModel = modelService.getApplication();
       
         switch(fieldName) {
            case "flowURI":
                var list_1 = [];
                var flows_1 = [];
                if (applicationModel) {
                    var triggerMappings = applicationModel.getTriggerFlowModelMaps();
                    triggerMappings.map(function(triggerMapping) {
                        console.log("trigger name = " + triggerMapping.getTriggerElement().name)
                        console.log("trigger flow = " + triggerMapping.getFlowModel().getName())
                        console.log("isTriggerFlow = " + triggerMapping.getFlowModel().isTriggerFlow())
                        console.log("current flow= " + context.getCurrentFlowName())

                        if (triggerMapping.getTriggerElement().name != "funit" &&
                            triggerMapping.getTriggerElement().name != "mock" &&
                                (context.getCurrentFlowName() !== triggerMapping.getFlowModel().getName()) && 
                                !triggerMapping.getFlowModel().isTriggerFlow()) {

                            if (list_1.indexOf(triggerMapping.getFlowModel().getName()) === -1) {
                                list_1.push(triggerMapping.getFlowModel().getName());
                                flows_1.push({
                                    "unique_id": "res://flow:" + triggerMapping.getFlowModel().getName().replace(/ /g, "_"),
                                    "name": triggerMapping.getFlowModel().getName()
                                });
                            }
                        }
                    });
                }
                return flows_1;
        }
        return null
    }

    validate = (fieldName: string, context: ITriggerContribution): Observable<IValidationResult> | IValidationResult => {
        return null
    }

    action = (actionId: string, context: ICreateFlowActionContext): Observable<IActionResult> | IActionResult => {
        let result = CreateFlowActionResult.newActionResult();
        return Observable.create(observer => {
            this.createFlow(context,  result);                                    
            let actionResult = ActionResult.newActionResult().setSuccess(true).setResult(result);
            observer.next(actionResult);
        });
    }

    createFlow(context, result) {
        let modelService = this.getModelService();
        var applicationModel = modelService.getApplication();
        var flow = context.getField("flowURI").value;
        
        let initrigger = modelService.createTriggerElement("funit/mock");
        if (initrigger) {
            for (let s = 0; s < initrigger.handler.settings.length; s++) {
                if (initrigger.handler.settings[s].name === "flowURI") {
                    initrigger.handler.settings[s].value = flow;
                } 
            }
           
        }

        let flowName = context.getFlowName();
        let iniflowModel = modelService.createFlow(flowName, context.getFlowDescription());
    /*    
        let inputschema = null
        let outputschema = null
        if (applicationModel && flow) {
            var triggerMappings = applicationModel.getTriggerFlowModelMaps();
            for (var i = 0; i < triggerMappings.length; i++) {
                var triggerMapping = triggerMappings[i];
                if (flow === "res://flow:" + triggerMapping.getFlowModel().getName().replace(/ /g, "_")) {
                    inputschema = triggerMapping.getFlowModel().getFlowInputSchema();
                    outputschema = triggerMapping.getFlowModel().getFlowOutputSchema();
                    break;
                }
            }
        }
        iniflowModel.addFlowInputSchema(inputschema);
        iniflowModel.addFlowOutputSchema(outputschema)
        */
        result = result.addTriggerFlowMapping(lodash.cloneDeep(initrigger), lodash.cloneDeep(iniflowModel));
    }
}