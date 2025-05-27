package myproject;

import com.pulumi.Context;
import com.pulumi.Pulumi;

import com.pulumi.kubernetes.apiextensions.v1.CustomResourceDefinition;
import com.pulumi.kubernetes.apiextensions.v1.CustomResourceDefinitionArgs;
import com.pulumi.kubernetes.apiextensions.v1.inputs.CustomResourceDefinitionNamesArgs;
import com.pulumi.kubernetes.apiextensions.v1.inputs.CustomResourceDefinitionVersionArgs;
import com.pulumi.kubernetes.apiextensions.v1.inputs.CustomResourceDefinitionSpecArgs;
import com.pulumi.kubernetes.apiextensions.v1.inputs.CustomResourceValidationArgs;
import com.pulumi.kubernetes.meta.v1.inputs.ObjectMetaArgs;
import com.pulumi.core.Output;
import com.pulumi.kubernetes.apiextensions.v1.inputs.JSONSchemaPropsArgs;
import java.util.HashMap;

public class App {
    public static void main(String[] args) {
        Pulumi.run(App::stack);
    }

    public static void stack(Context ctx) {
        var metadata = ObjectMetaArgs.builder()
                .name("javacrds.example.com")
                .build();

        var spec = CustomResourceDefinitionSpecArgs.builder()
                .group("example.com")
                .scope("Namespaced")
                .names(CustomResourceDefinitionNamesArgs.builder()
                        .kind("JavaCRD")
                        .plural("javacrds")
                        .singular("javacrd")
                        .shortNames("jcrd")
                        .build())
                .versions(CustomResourceDefinitionVersionArgs.builder()
                        .name("v1")
                        .served(true)
                        .storage(true)
                        .schema(CustomResourceValidationArgs.builder()
                                .openAPIV3Schema(JSONSchemaPropsArgs.builder()
                                        .type("object")
                                        .properties(new HashMap<String, JSONSchemaPropsArgs>() {
                                            {
                                                put("key", JSONSchemaPropsArgs.builder()
                                                        .type("object")
                                                        .x_kubernetes_preserve_unknown_fields(true)
                                                        .build());
                                            }
                                        }).build())
                                .build())
                        .build())
                .build();

        new CustomResourceDefinition("crd",
                CustomResourceDefinitionArgs.builder()
                        .metadata(metadata)
                        .spec(spec)
                        .build());

        var crdGet = CustomResourceDefinition.get("getCRDUrn", Output.of("javacrds.example.com"), null);
        ctx.export("urn", crdGet.urn());
    }
}
