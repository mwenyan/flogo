import {Observable} from "rxjs/Observable";
import {Injectable, Injector, Inject} from "@angular/core";
import {Http} from "@angular/http";
import {
    WiContrib,
    WiServiceHandlerContribution,
    IValidationResult,
    ValidationResult,
    IActivityContribution,
    WiContribModelService,
} from "wi-studio/app/contrib/wi-contrib";


@WiContrib({})
@Injectable()
export class MockActivityContributionHandler extends WiServiceHandlerContribution {
    constructor(private injector: Injector, private http: Http, private contribModelService: WiContribModelService) {
        super(injector, http, contribModelService);
    }
   
    value = (fieldName: string, context: IActivityContribution): any | Observable<any> => {
        var modelService = this.getModelService();
        var applicationModel = modelService.getApplication();
        var flow = context.getField("flowURI").value;
        if (fieldName === "flowURI") {
            var list_1 = [];
            var flows_1 = [];
            if (applicationModel) {
                var triggerMappings = applicationModel.getTriggerFlowModelMaps();
                triggerMappings.map(function(triggerMapping) {
                    if ((context.getCurrentFlowName() !== triggerMapping.getFlowModel().getName()) && !triggerMapping.getFlowModel().isTriggerFlow()) {
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
        } else if (fieldName === "input") {
            if (applicationModel && flow) {
                var triggerMappings = applicationModel.getTriggerFlowModelMaps();
                var schema = null;
                for (var i = 0; i < triggerMappings.length; i++) {
                    var triggerMapping = triggerMappings[i];
                    if (flow === "res://flow:" + triggerMapping.getFlowModel().getName().replace(/ /g, "_")) {
                        schema = triggerMapping.getFlowModel().getFlowInputSchema().json;
                        break;
                    }
                }
                return schema;
            }
        } else if (fieldName === "output") {
            if (applicationModel && flow) {
                var triggerMappings = applicationModel.getTriggerFlowModelMaps();
                var schema = null;
                for (var i = 0; i < triggerMappings.length; i++) {
                    var triggerMapping = triggerMappings[i];
                    if (flow === "res://flow:" + triggerMapping.getFlowModel().getName().replace(/ /g, "_")) {
                        schema = triggerMapping.getFlowModel().getFlowOutputSchema().json;
                        break;
                    }
                }
                return schema;
            }
        }
        return null;
    }

    validate = (fieldName: string, context: IActivityContribution): Observable<IValidationResult> | IValidationResult => {
        var modelService = this.getModelService();
        var applicationModel = modelService.getApplication();
        var flow = context.getField("flowURI").value;
        if (fieldName === "flowURI" && flow) {
            var list_2 = [];
            var flows_2 = {};
            if (applicationModel) {
                var triggerMappings = applicationModel.getTriggerFlowModelMaps();
                triggerMappings.map(function(triggerMapping) {
                    if (!triggerMapping.getFlowModel().isTriggerFlow()) {
                        var name_1 = "res://flow:" + triggerMapping.getFlowModel().getName().replace(/ /g, "_");
                        if (list_2.indexOf(name_1) === -1) {
                            list_2.push(name_1);
                            flows_2[name_1] = triggerMapping.getFlowModel();
                        }
                    }
                });
                return Observable.create(observer => {
                    let vresult: IValidationResult = ValidationResult.newValidationResult();
                    if (this.isLoopDetected(flow, flow, flows_2)) {
                        vresult.setError("SUBFLOW-LOOP", "Cyclic dependency detected in the subflow");
                    } 
                   
                    observer.next(vresult);
                });
            }
        }
        return null;
    }

    isLoopDetected = function(destination, source, flows) {
        if (flows && flows[source]) {
            var errorFlow = flows[source].getErrorFlow();
            if (this.checkLoopInFlow(destination, source, flows[source], flows)) {
                return true;
            } else if (errorFlow && this.checkLoopInFlow(destination, source, errorFlow, flows)) {
                return true;
            }
        }
        return false;
    }
    
    checkLoopInFlow = function(destination, source, flow, flows) {
        var element = flow.getFlowElement();
        if (this.checkCurrentElement(element, destination, flows) || this.checkChildren(element, destination, flows)) {
            return true;
        }
        return false;
    }
    ;
    checkChildren = function(element, destination, flows) {
        var children = element.getChildren();
        var isError = false;
        for (var _i = 0, children_1 = children; _i < children_1.length; _i++) {
            var child = children_1[_i];
            if (this.checkCurrentElement(child, destination, flows)) {
                isError = true;
                break;
            } else {
                isError = this.checkChildren(child, destination, flows);
            }
        }
        return isError;
    }
    
    checkCurrentElement = function(element, destination, flows) {
        if (element.name === "flogo-subflow") {
            var flow = element.getField("flowURI").value;
            if (flow === destination) {
                return true;
            } else {
                return this.isLoopDetected(destination, flow, flows);
            }
        }
        return false;
    }
}