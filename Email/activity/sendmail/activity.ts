import { Http } from "@angular/http";
import { Injectable, Inject, Injector } from "@angular/core";
import {
  WiContrib,
  WiServiceProviderContribution,
  AbstractContribFieldProvider,
  AbstractContribValidationProvider
} from "wi-studio/app/contrib/wi-contrib";
import { IActivityContribution, IContribValidationProvider, IContributionContext } from "wi-studio/common/models/contrib";
import { IValidationResult, ValidationResult } from "wi-studio/common/models/validation";
import { Observable } from "rxjs/Observable";

@Injectable()
export class UserNameFieldProvider extends AbstractContribValidationProvider {
  constructor() {
    super();
  }
  validate(context: IActivityContribution): Observable<IValidationResult> {
    return Observable.create(observer => {
      for (let i = 0; i < context.inputs.length; i++) {
        if (context.inputs[i].name.toLowerCase() === "connection security") {
          if (context.inputs[i].value && context.inputs[i].value.length > 0 && context.inputs[i].value.toLowerCase() === "none") {
            observer.next(ValidationResult.newValidationResult().setVisible(false));
            break;
          } else {
            observer.next(ValidationResult.newValidationResult().setVisible(true));
            break;
          }
        }
      }
      observer.complete();
    });
  }
}

@Injectable()
export class PasswordFieldProvider extends AbstractContribValidationProvider {
  constructor() {
    super();
  }
  validate(context: IActivityContribution): Observable<IValidationResult> {
    return Observable.create(observer => {
      for (let i = 0; i < context.inputs.length; i++) {
        if (context.inputs[i].name.toLowerCase() === "connection security") {
          if (context.inputs[i].value && context.inputs[i].value.length > 0 && context.inputs[i].value.toLowerCase() === "none") {
            observer.next(ValidationResult.newValidationResult().setVisible(false));
            break;
          } else {
            observer.next(ValidationResult.newValidationResult().setVisible(true));
            break;
          }
        }
      }
      observer.complete();
    });
  }
}

@Injectable()
@WiContrib({
  validationProviders: [
    {
      field: "Username",
      useClass: UserNameFieldProvider
    },
    {
      field: "Password",
      useClass: PasswordFieldProvider
    }
  ],
  fieldProviders: [

  ]
})
export class SendMailActivityService extends WiServiceProviderContribution {
  constructor( @Inject(Injector) injector, @Inject(Http) http: Http) {
    super(injector, http);
  }
}

