package myproject;

import java.util.Map;
import java.util.Objects;
import java.util.Optional;

import javax.annotation.Nullable;

import com.pulumi.Pulumi;
import com.pulumi.core.Output;
import com.pulumi.core.annotations.CustomType;
import com.pulumi.core.annotations.Export;
import com.pulumi.core.annotations.Import;
import com.pulumi.kubernetes.apiextensions.CustomResource;
import com.pulumi.kubernetes.apiextensions.CustomResourceArgs;
import com.pulumi.kubernetes.apiextensions.CustomResourceArgsBase;
import com.pulumi.kubernetes.meta.v1.inputs.ObjectMetaArgs;

/**
 * Demonstration of working with custom resources in the Java SDK.
 * Prerequisites:
 * - cert-manager v1.x installed in your Kubernetes cluster.
 * 
 * This example deploys a cert-manager.io/v1 Issuer to your Kubernetes cluster.
 * Two different ways of defining the Issuer are demonstrated:
 * - Using the CustomResource class directly.
 * - Using a custom Issuer class that extends CustomResource, to provide a
 * type-safe API.
 */
public class App {
    public static void main(String[] args) {
        Pulumi.run(ctx -> {
            var issuer1 = new CustomResource("issuer1", CustomResourceArgs.builder()
                    .apiVersion("cert-manager.io/v1")
                    .kind("Issuer")
                    .metadata(ObjectMetaArgs.builder().build())
                    .build());
            ctx.export("issuer1_name", issuer1.metadata().applyValue(s -> s.name()));

            var issuer2 = new Issuer("issuer2", IssuerArgs.builder()
                    .metadata(ObjectMetaArgs.builder().build())
                    .spec(Inputs.IssuerSpecArgs.builder()
                            .selfSigned(Inputs.SelfSignedArgs.builder().build())
                            .build())
                    .build());

            ctx.export("issuer2_name", issuer2.metadata().applyValue(s -> s.name()));
            ctx.export("issuer2_selfsigned", issuer2.spec().applyValue(s -> s.selfSigned().isPresent()));
        });
    }
}

class Issuer extends CustomResource {
    /**
     * The spec of the Issuer.
     */
    @Export(name = "spec", refs = { Outputs.IssuerSpec.class }, tree = "[0]")
    private Output<Outputs.IssuerSpec> spec;

    public Output<Outputs.IssuerSpec> spec() {
        return this.spec;
    }

    public Issuer(String name, @Nullable IssuerArgs args) {
        super(name, makeArgs(args));
    }

    public Issuer(String name, @Nullable IssuerArgs args,
            @Nullable com.pulumi.resources.CustomResourceOptions options) {
        super(name, makeArgs(args), options);
    }

    protected Issuer(String name, Output<String> id,
            @Nullable com.pulumi.resources.CustomResourceOptions options) {
        super(name, "cert-manager.io/v1", "Issuer", id, options);
    }

    private static IssuerArgs makeArgs(@Nullable IssuerArgs args) {
        var builder = args == null ? IssuerArgs.builder() : IssuerArgs.builder(args);
        return builder
                .apiVersion("cert-manager.io/v1")
                .kind("Issuer")
                .build();
    }

    public static Issuer get(String name, Output<String> id,
            @Nullable com.pulumi.resources.CustomResourceOptions options) {
        return new Issuer(name, id, options);
    }
}

class IssuerArgs extends CustomResourceArgsBase {
    /**
     * The spec of the Issuer.
     */
    @Import(name = "spec", required = true)
    @Nullable
    private Output<Inputs.IssuerSpecArgs> spec;

    public static Builder builder() {
        return new Builder();
    }

    public static Builder builder(IssuerArgs defaults) {
        return new Builder(defaults);
    }

    static class Builder extends CustomResourceArgsBase.Builder<IssuerArgs, Builder> {
        public Builder() {
            super(new IssuerArgs());
        }

        public Builder(IssuerArgs defaults) {
            super(new IssuerArgs(), defaults);
        }

        public Builder spec(@Nullable Output<Inputs.IssuerSpecArgs> spec) {
            $.spec = spec;
            return this;
        }

        public Builder spec(Inputs.IssuerSpecArgs spec) {
            return spec(Output.of(spec));
        }

        @Override
        protected void copy(IssuerArgs args) {
            super.copy(args);
            $.spec = args.spec;
        }
    }
}

class Inputs {
    public static final class IssuerSpecArgs extends com.pulumi.resources.ResourceArgs {

        public static final IssuerSpecArgs Empty = new IssuerSpecArgs();

        @Import(name = "selfSigned")
        private @Nullable Output<SelfSignedArgs> selfSigned;

        public Optional<Output<SelfSignedArgs>> selfSigned() {
            return Optional.ofNullable(this.selfSigned);
        }

        private IssuerSpecArgs() {
        }

        private IssuerSpecArgs(IssuerSpecArgs $) {
            this.selfSigned = $.selfSigned;
        }

        public static Builder builder() {
            return new Builder();
        }

        public static Builder builder(IssuerSpecArgs defaults) {
            return new Builder(defaults);
        }

        public static final class Builder {
            private IssuerSpecArgs $;

            public Builder() {
                $ = new IssuerSpecArgs();
            }

            public Builder(IssuerSpecArgs defaults) {
                $ = new IssuerSpecArgs(Objects.requireNonNull(defaults));
            }

            public Builder selfSigned(@Nullable Output<SelfSignedArgs> selfSigned) {
                $.selfSigned = selfSigned;
                return this;
            }

            public Builder selfSigned(SelfSignedArgs selfSigned) {
                return selfSigned(Output.of(selfSigned));
            }

            public IssuerSpecArgs build() {
                return $;
            }
        }
    }

    public static final class SelfSignedArgs extends com.pulumi.resources.ResourceArgs {

        public static final SelfSignedArgs Empty = new SelfSignedArgs();

        private SelfSignedArgs() {
        }

        private SelfSignedArgs(SelfSignedArgs $) {
        }

        public static Builder builder() {
            return new Builder();
        }

        public static Builder builder(SelfSignedArgs defaults) {
            return new Builder(defaults);
        }

        public static final class Builder {
            private SelfSignedArgs $;

            public Builder() {
                $ = new SelfSignedArgs();
            }

            public Builder(SelfSignedArgs defaults) {
                $ = new SelfSignedArgs(Objects.requireNonNull(defaults));
            }

            public SelfSignedArgs build() {
                return $;
            }
        }
    }
}

class Outputs {
    @CustomType
    static final class IssuerSpec {

        private @Nullable SelfSigned selfSigned;

        private IssuerSpec() {
        }

        public Optional<SelfSigned> selfSigned() {
            return Optional.ofNullable(this.selfSigned);
        }

        public static Builder builder() {
            return new Builder();
        }

        public static Builder builder(IssuerSpec defaults) {
            return new Builder(defaults);
        }

        @CustomType.Builder
        public static final class Builder {
            private @Nullable SelfSigned selfSigned;

            public Builder() {
            }

            public Builder(IssuerSpec defaults) {
                Objects.requireNonNull(defaults);
                this.selfSigned = defaults.selfSigned;
            }

            @CustomType.Setter
            public Builder selfSigned(@Nullable SelfSigned selfSigned) {
                this.selfSigned = selfSigned;
                return this;
            }

            public IssuerSpec build() {
                final var _resultValue = new IssuerSpec();
                _resultValue.selfSigned = selfSigned;
                return _resultValue;
            }
        }
    }

    @CustomType
    static final class SelfSigned {

        private SelfSigned() {
        }

        public static Builder builder() {
            return new Builder();
        }

        public static Builder builder(SelfSigned defaults) {
            return new Builder(defaults);
        }

        @CustomType.Builder
        public static final class Builder {
            public Builder() {
            }

            public Builder(SelfSigned defaults) {
                Objects.requireNonNull(defaults);
            }

            public SelfSigned build() {
                final var _resultValue = new SelfSigned();
                return _resultValue;
            }
        }
    }
}
