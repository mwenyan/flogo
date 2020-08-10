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
export class ListDocumentActivityContributionHandler extends WiServiceHandlerContribution {
    constructor(private injector: Injector, private http: Http,) {
        super(injector, http);
    }
   
    value = (fieldName: string, context: IActivityContribution): any | Observable<any> => {

        switch(fieldName) {
            case "connector":
                let connections = [];
                return Observable.create(observer => {
                    WiContributionUtils.getConnections(this.http, "AzureCosmosDB", "CosmosDBConnector").subscribe((data: IConnectorContribution[]) => {
                        data.forEach(connection => {
                            if ((<any>connection).isValid) {
                                for(let i=0; i < connection.settings.length; i++) {
                                    if(connection.settings[i].name === "name"){
                                        connections.push({
                                            "unique_id": WiContributionUtils.getUniqueId(connection),
                                            "name": connection.settings[i].value
                                        });
                                        break;
                                    }
                                }
                            }
                        });
                        observer.next(connections);
                    });
                });

            case "input":
                let schema = {
                    "$schema": "http://json-schema.org/draft-04/schema#",
                    "type": "object",
                    "properties": {
                        "database": {
                            "type": "string"
                        },
                        "collection": {
                            "type": "string"
                        },
                        "x-ms-activity-id": {
                            "type": "string"
                        },
                        "x-ms-max-item-count": {
                            "type": "integer"
                        },
                        "x-ms-continuation": {
                            "type": "string"
                        },
                        "x-ms-consistency-level": {
                            "type": "string"
                        },
                        "x-ms-session-token": {
                            "type": "string"
                        },
                        "If-None-Match": {
                            "type": "string"
                        },
                        "x-ms-documentdb-partitionkeyrangeid": {
                            "type": "string"
                        }
                    },
                    "required":["database", "collection"]
                }

                return JSON.stringify(schema);;
            case "output":
                let oschemaobj = {code: "", message: "", "x-ms-activity-id": "", "x-ms-session-token":"", "x-ms-continuation": "", "_rid": "", "Documents":[], "_count": 0}
               
                let oschema = {
                    "$schema": "http://json-schema.org/draft-04/schema#",
                    "type": "object",
                    "properties": {
                        "code": {
                            "type": "string"
                        },
                        "message": {
                            "type": "string"
                        },
                        "x-ms-activity-id": {
                            "type": "string"
                        },
                        "x-ms-session-token": {
                            "type": "string"
                        },
                        "x-ms-continuation": {
                            "type": "string"
                        },
                        "_rid": {
                            "type": "string"
                        },
                        "Documents": {
                            "type": "array",
                            "items": {
                                "type": "object",
                                "properties": {
                                }
                            }
                        },
                        "_count": {
                            "type": "integer"
                        },
                    }
                }
                let doc = context.getField("document").value;
                if(Boolean(doc)){
                    let docschema = JSON.parse(doc);
                    if(Boolean(docschema["$schema"])){
                        docschema["properties"]["id"] = {type: "string"}
                        docschema["properties"]["_rid"] = {type: "string"}
                        docschema["properties"]["_self"] = {type: "string"}
                        docschema["properties"]["_etag"] = {type: "string"}
                        docschema["properties"]["_attachments"] = {type: "string"}
                        docschema["properties"]["_ts"] = {type: "integer"}
                        oschema.properties.Documents.items.properties = docschema["properties"]
                           
                        return JSON.stringify(oschema)
                    } else {
                        docschema["id"] = ""
                        docschema["_rid"] = ""
                        docschema["_self"] = ""
                        docschema["_etag"] = ""
                        docschema["_attachments"] = ""
                        docschema["_ts"] = 0
                        oschemaobj.Documents.push(docschema)
                           
                        return JSON.stringify(oschemaobj)
                    }
                }

                return null;
                
        }
        return null;
    }
 
    validate = (fieldName: string, context: IActivityContribution): Observable<IValidationResult> | IValidationResult => {
        return null;
    }
/*
    listDatabases(conId):  Observable<any> {
        return Observable.create(observer => {
            WiContributionUtils.getConnection(this.http, conId)
                            .map(data => data)
                            .subscribe(data => {
                                let account, masterkey;
                                for (let setting of data.settings) {
                                    if(setting.name === "account") {
                                        account = setting.value;
                                        break;
                                    } else if(setting.name === "key") {
                                        masterkey = setting.value;
                                        break;
                                    }
                                }
                                let today = this.getToday();
                                    // assign our verb
                                var verb = "get";
                                var resType = "dbs";
                                var resourceId = "dbs";
                                let auth = this.getAuthToken( masterkey, today, verb, resType, resourceId)
                                const headers = { 'Authorization': auth, 'Accept': 'application/json', 'x-ms-version': '2016-07-11', 'x-ms-date': today}
                                const body = {}
                                this.http.post<any>('https://' + this.getHostName(account) + '/dbs', body, { headers }).subscribe(data => {
                                    let dbs = []   
                                    data.Databases.array.forEach(element => {
                                        dbs.push({
                                            "unique_id": WiContributionUtils.getUniqueId(element.id),
                                            "name": element.id
                                        });
                                    }); 
                                    return observer.next(dbs);
                                })
                            });
                        });
    }

    getHostName(account): string {
        return account + ".documents.azure.com:443"
    }

    getToday(): string {
        var today = new Date();
        var UTCstring = today.toUTCString();
        return UTCstring.toLowerCase();
    }
    getAuthToken(masterkey, date, verb, resType, resourceId): string {
        // parse our master key out as base64 encoding
        var key = CryptoJS.enc.Base64.parse(masterkey);
       
        // build up the request text for the signature so can sign it along with the key
        var text = (verb || "").toLowerCase() + "\n" + 
                    (resType || "").toLowerCase() + "\n" + 
                    (resourceId || "") + "\n" + 
                    (date || "").toLowerCase() + "\n" + 
                    "" + "\n";
        
        // create the signature from build up request text
        var signature = CryptoJS.HmacSHA256(text, key);
        
        // back to base 64 bits
        var base64Bits = CryptoJS.enc.Base64.stringify(signature);
        
        // format our authentication token and URI encode it.
        var MasterToken = "master";
        var TokenVersion = "1.0";
        var auth = encodeURIComponent("type=" + MasterToken + "&ver=" + TokenVersion + "&sig=" + base64Bits);

        return auth;
    }

    getRequestUrl(account, database, collection) {
        return "https://" + account + "/"
    }

    listCollections(conId, database):  Observable<any> {

        return Observable.create(observer => {
            WiContributionUtils.getConnection(this.http, conId)
                            .map(data => data)
                            .subscribe(data => {
                                let account, masterkey;
                                for (let setting of data.settings) {
                                    if(setting.name === "account") {
                                        account = setting.value;
                                        break;
                                    } else if(setting.name === "key") {
                                        masterkey = setting.value;
                                        break;
                                    }
                                }
                                let today = this.getToday();
                                    // assign our verb
                                var verb = "get";
                                var resType = "colls";
                                var resourceId = "colls";
                                let auth = this.getAuthToken(masterkey, today, verb, resType, resourceId);
                                const headers = { 'Authorization': auth, 'Accept': 'application/json', 'x-ms-version': '2016-07-11', 'x-ms-date': today}
                                const body = {}
                                this.http.post<any>('https://' + this.getHostName(account) + '/dbs/' + database + '/colls', body, { headers }).subscribe(data => {
                                    let colls = []   
                                    data.DocumentCollections.array.forEach(element => {
                                        colls.push({
                                            "unique_id": WiContributionUtils.getUniqueId(element.id),
                                            "name": element.id
                                        });
                                    }); 
                                    return observer.next(colls);
                                })
                            });
                        });
    }
    */
}