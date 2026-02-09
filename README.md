# Testing the Identity Hub
An error occured during one test run of the local docker compose setup. 
Since it has occured once again, this repository contains a test to find out if there is actually something wrong. 

It will create a whole lot of requests and it will log out runs that failed to logs/failed as well as some random runs to ok for comparison.

To run, you need Go installed (>= 1.25 at least). Then run:

~~~
go run .
~~~

from the project root folder. 

## Structure of jsons
The produced jsons in failed contain: 

- Did: DID used for request
- ParticipantRequest: The actual data that was sent (not encoded as it is from the struct directly)
- ParticipantRequestResponse: What the server responded with
- DidDocumentResponse: Just the did document of this run
- KeyPairResponse: left empty, but the [IdentityHub API](https://eclipse-edc.github.io/IdentityHub/openapi/identity-api/#/) has a KeyPair endpoint to get all of them - just remember to put ?limit=10000 to get all.

# Results
So far, I was able to reproduce the error about 20 times (give or take 10) for every 5000 runs.

What I was able to see: 
None of them seem out of place at all. Apart from the fact, that all failed request have an "X" Property within their according keypairs.json entry (identified using just test-id-xxxx) which has one char less then the others. But that could simply be a symptom instead of the actual cause.

The logs of consumer_id_hub always show the same error - and a lot of warnings about some state? But that is always there, even on good ones.

This is the error popping up: 
~~~log
java.security.spec.InvalidKeySpecException: java.security.InvalidKeyException: key length must be 32
	at java.base/sun.security.ec.ed.EdDSAKeyFactory.engineGeneratePublic(Unknown Source)
	at java.base/java.security.KeyFactory.generatePublic(Unknown Source)
	at org.eclipse.edc.keys.keyparsers.JwkParser.readPublicKey(JwkParser.java:186)
	at org.eclipse.edc.keys.keyparsers.JwkParser.parseOctetKeyPair(JwkParser.java:143)
	at org.eclipse.edc.keys.keyparsers.JwkParser.parse(JwkParser.java:94)
	at org.eclipse.edc.keys.KeyParserRegistryImpl.lambda$parse$1(KeyParserRegistryImpl.java:38)
	at java.base/java.util.Optional.map(Unknown Source)
	at org.eclipse.edc.keys.KeyParserRegistryImpl.parse(KeyParserRegistryImpl.java:38)
	at org.eclipse.edc.identityhub.did.DidDocumentServiceImpl.lambda$keyPairActivated$16(DidDocumentServiceImpl.java:267)
	at org.eclipse.edc.transaction.local.LocalTransactionContext.lambda$execute$0(LocalTransactionContext.java:57)
	at org.eclipse.edc.transaction.local.LocalTransactionContext.execute(LocalTransactionContext.java:74)
	at org.eclipse.edc.transaction.local.LocalTransactionContext.execute(LocalTransactionContext.java:56)
	at org.eclipse.edc.identityhub.did.DidDocumentServiceImpl.keyPairActivated(DidDocumentServiceImpl.java:259)
	at org.eclipse.edc.identityhub.did.DidDocumentServiceImpl.on(DidDocumentServiceImpl.java:252)
	at org.eclipse.edc.runtime.core.event.EventRouterImpl.lambda$publish$2(EventRouterImpl.java:73)
	at java.base/java.util.stream.ForEachOps$ForEachOp$OfRef.accept(Unknown Source)
	at java.base/java.util.ArrayList$ArrayListSpliterator.forEachRemaining(Unknown Source)
	at java.base/java.util.stream.ReferencePipeline$Head.forEach(Unknown Source)
	at java.base/java.util.stream.ReferencePipeline$7$1FlatMap.accept(Unknown Source)
	at java.base/java.util.stream.ReferencePipeline$2$1.accept(Unknown Source)
	at java.base/java.util.concurrent.ConcurrentHashMap$EntrySpliterator.forEachRemaining(Unknown Source)
	at java.base/java.util.stream.AbstractPipeline.copyInto(Unknown Source)
	at java.base/java.util.stream.AbstractPipeline.wrapAndCopyInto(Unknown Source)
	at java.base/java.util.stream.ForEachOps$ForEachOp.evaluateSequential(Unknown Source)
	at java.base/java.util.stream.ForEachOps$ForEachOp$OfRef.evaluateSequential(Unknown Source)
	at java.base/java.util.stream.AbstractPipeline.evaluate(Unknown Source)
	at java.base/java.util.stream.ReferencePipeline.forEach(Unknown Source)
	at org.eclipse.edc.runtime.core.event.EventRouterImpl.publish(EventRouterImpl.java:73)
	at org.eclipse.edc.identityhub.keypairs.KeyPairEventPublisher.publish(KeyPairEventPublisher.java:89)
	at org.eclipse.edc.identityhub.keypairs.KeyPairEventPublisher.activated(KeyPairEventPublisher.java:81)
	at org.eclipse.edc.identityhub.keypairs.KeyPairServiceImpl.lambda$activateKeyPair$17(KeyPairServiceImpl.java:235)
	at java.base/java.util.concurrent.ConcurrentLinkedQueue.forEachFrom(Unknown Source)
	at java.base/java.util.concurrent.ConcurrentLinkedQueue.forEach(Unknown Source)
	at org.eclipse.edc.spi.observe.Observable.invokeForEach(Observable.java:47)
	at org.eclipse.edc.identityhub.keypairs.KeyPairServiceImpl.lambda$activateKeyPair$18(KeyPairServiceImpl.java:235)
	at org.eclipse.edc.spi.result.AbstractResult.onSuccess(AbstractResult.java:82)
	at org.eclipse.edc.identityhub.keypairs.KeyPairServiceImpl.activateKeyPair(KeyPairServiceImpl.java:235)
	at org.eclipse.edc.identityhub.keypairs.KeyPairServiceImpl.lambda$addKeyPair$4(KeyPairServiceImpl.java:119)
	at org.eclipse.edc.spi.result.AbstractResult.compose(AbstractResult.java:185)
	at org.eclipse.edc.identityhub.keypairs.KeyPairServiceImpl.lambda$addKeyPair$5(KeyPairServiceImpl.java:117)
	at org.eclipse.edc.transaction.local.LocalTransactionContext.execute(LocalTransactionContext.java:74)
	at org.eclipse.edc.identityhub.keypairs.KeyPairServiceImpl.addKeyPair(KeyPairServiceImpl.java:76)
	at org.eclipse.edc.identityhub.participantcontext.ParticipantContextEventCoordinator.lambda$on$0(ParticipantContextEventCoordinator.java:79)
	at org.eclipse.edc.spi.result.AbstractResult.compose(AbstractResult.java:185)
	at org.eclipse.edc.identityhub.participantcontext.ParticipantContextEventCoordinator.on(ParticipantContextEventCoordinator.java:79)
	at org.eclipse.edc.runtime.core.event.EventRouterImpl.lambda$publish$2(EventRouterImpl.java:73)
	at java.base/java.util.stream.ForEachOps$ForEachOp$OfRef.accept(Unknown Source)
	at java.base/java.util.ArrayList$ArrayListSpliterator.forEachRemaining(Unknown Source)
	at java.base/java.util.stream.ReferencePipeline$Head.forEach(Unknown Source)
	at java.base/java.util.stream.ReferencePipeline$7$1FlatMap.accept(Unknown Source)
	at java.base/java.util.stream.ReferencePipeline$2$1.accept(Unknown Source)
	at java.base/java.util.concurrent.ConcurrentHashMap$EntrySpliterator.forEachRemaining(Unknown Source)
	at java.base/java.util.stream.AbstractPipeline.copyInto(Unknown Source)
	at java.base/java.util.stream.AbstractPipeline.wrapAndCopyInto(Unknown Source)
	at java.base/java.util.stream.ForEachOps$ForEachOp.evaluateSequential(Unknown Source)
	at java.base/java.util.stream.ForEachOps$ForEachOp$OfRef.evaluateSequential(Unknown Source)
	at java.base/java.util.stream.AbstractPipeline.evaluate(Unknown Source)
	at java.base/java.util.stream.ReferencePipeline.forEach(Unknown Source)
	at org.eclipse.edc.runtime.core.event.EventRouterImpl.publish(EventRouterImpl.java:73)
	at org.eclipse.edc.identityhub.participantcontext.ParticipantContextEventPublisher.publish(ParticipantContextEventPublisher.java:79)
	at org.eclipse.edc.identityhub.participantcontext.ParticipantContextEventPublisher.created(ParticipantContextEventPublisher.java:45)
	at org.eclipse.edc.identityhub.participantcontext.ParticipantContextServiceImpl.lambda$createParticipantContext$2(ParticipantContextServiceImpl.java:90)
	at java.base/java.util.concurrent.ConcurrentLinkedQueue.forEachFrom(Unknown Source)
	at java.base/java.util.concurrent.ConcurrentLinkedQueue.forEach(Unknown Source)
	at org.eclipse.edc.spi.observe.Observable.invokeForEach(Observable.java:47)
	at org.eclipse.edc.identityhub.participantcontext.ParticipantContextServiceImpl.lambda$createParticipantContext$3(ParticipantContextServiceImpl.java:90)
	at org.eclipse.edc.spi.result.AbstractResult.onSuccess(AbstractResult.java:82)
	at org.eclipse.edc.identityhub.participantcontext.ParticipantContextServiceImpl.lambda$createParticipantContext$4(ParticipantContextServiceImpl.java:90)
	at org.eclipse.edc.transaction.local.LocalTransactionContext.execute(LocalTransactionContext.java:74)
	at org.eclipse.edc.identityhub.participantcontext.ParticipantContextServiceImpl.createParticipantContext(ParticipantContextServiceImpl.java:74)
	at org.eclipse.edc.identityhub.api.verifiablecredential.v1.unstable.ParticipantContextApiController.createParticipant(ParticipantContextApiController.java:69)
	at java.base/jdk.internal.reflect.DirectMethodHandleAccessor.invoke(Unknown Source)
	at java.base/java.lang.reflect.Method.invoke(Unknown Source)
	at org.glassfish.jersey.server.model.internal.ResourceMethodInvocationHandlerFactory.lambda$static$0(ResourceMethodInvocationHandlerFactory.java:52)
	at org.glassfish.jersey.server.model.internal.AbstractJavaResourceMethodDispatcher$1.run(AbstractJavaResourceMethodDispatcher.java:146)
	at org.glassfish.jersey.server.model.internal.AbstractJavaResourceMethodDispatcher.invoke(AbstractJavaResourceMethodDispatcher.java:189)
	at org.glassfish.jersey.server.model.internal.JavaResourceMethodDispatcherProvider$TypeOutInvoker.doDispatch(JavaResourceMethodDispatcherProvider.java:219)
	at org.glassfish.jersey.server.model.internal.AbstractJavaResourceMethodDispatcher.dispatch(AbstractJavaResourceMethodDispatcher.java:93)
	at org.glassfish.jersey.server.model.ResourceMethodInvoker.invoke(ResourceMethodInvoker.java:478)
	at org.glassfish.jersey.server.model.ResourceMethodInvoker.apply(ResourceMethodInvoker.java:400)
	at org.glassfish.jersey.server.model.ResourceMethodInvoker.apply(ResourceMethodInvoker.java:81)
	at org.glassfish.jersey.server.ServerRuntime$1.run(ServerRuntime.java:274)
	at org.glassfish.jersey.internal.Errors$1.call(Errors.java:248)
	at org.glassfish.jersey.internal.Errors$1.call(Errors.java:244)
	at org.glassfish.jersey.internal.Errors.process(Errors.java:292)
	at org.glassfish.jersey.internal.Errors.process(Errors.java:274)
	at org.glassfish.jersey.internal.Errors.process(Errors.java:244)
	at org.glassfish.jersey.process.internal.RequestScope.runInScope(RequestScope.java:266)
	at org.glassfish.jersey.server.ServerRuntime.process(ServerRuntime.java:253)
	at org.glassfish.jersey.server.ApplicationHandler.handle(ApplicationHandler.java:696)
	at org.glassfish.jersey.servlet.WebComponent.serviceImpl(WebComponent.java:397)
	at org.glassfish.jersey.servlet.WebComponent.service(WebComponent.java:349)
	at org.glassfish.jersey.servlet.ServletContainer.service(ServletContainer.java:358)
	at org.glassfish.jersey.servlet.ServletContainer.service(ServletContainer.java:312)
	at org.glassfish.jersey.servlet.ServletContainer.service(ServletContainer.java:205)
	at org.eclipse.jetty.ee10.servlet.ServletHolder.handle(ServletHolder.java:736)
	at org.eclipse.jetty.ee10.servlet.ServletHandler$ChainEnd.doFilter(ServletHandler.java:1622)
	at org.eclipse.jetty.ee10.servlet.ServletHandler$MappedServlet.handle(ServletHandler.java:1555)
	at org.eclipse.jetty.ee10.servlet.ServletChannel.dispatch(ServletChannel.java:823)
	at org.eclipse.jetty.ee10.servlet.ServletChannel.handle(ServletChannel.java:440)
	at org.eclipse.jetty.ee10.servlet.ServletHandler.handle(ServletHandler.java:470)
	at org.eclipse.jetty.server.handler.ContextHandler.handle(ContextHandler.java:1071)
	at org.eclipse.jetty.server.handler.ContextHandlerCollection.handle(ContextHandlerCollection.java:181)
	at org.eclipse.jetty.server.Server.handle(Server.java:182)
	at org.eclipse.jetty.server.internal.HttpChannelState$HandlerInvoker.run(HttpChannelState.java:678)
	at org.eclipse.jetty.server.internal.HttpConnection.onFillable(HttpConnection.java:416)
	at org.eclipse.jetty.io.AbstractConnection$ReadCallback.succeeded(AbstractConnection.java:322)
	at org.eclipse.jetty.io.FillInterest.fillable(FillInterest.java:99)
	at org.eclipse.jetty.io.SelectableChannelEndPoint$1.run(SelectableChannelEndPoint.java:53)
	at org.eclipse.jetty.util.thread.strategy.AdaptiveExecutionStrategy.runTask(AdaptiveExecutionStrategy.java:480)
	at org.eclipse.jetty.util.thread.strategy.AdaptiveExecutionStrategy.consumeTask(AdaptiveExecutionStrategy.java:443)
	at org.eclipse.jetty.util.thread.strategy.AdaptiveExecutionStrategy.tryProduce(AdaptiveExecutionStrategy.java:293)
	at org.eclipse.jetty.util.thread.strategy.AdaptiveExecutionStrategy.run(AdaptiveExecutionStrategy.java:201)
	at org.eclipse.jetty.util.thread.ReservedThreadExecutor$ReservedThread.run(ReservedThreadExecutor.java:311)
	at org.eclipse.jetty.util.thread.QueuedThreadPool.runJob(QueuedThreadPool.java:981)
	at org.eclipse.jetty.util.thread.QueuedThreadPool$Runner.doRunJob(QueuedThreadPool.java:1211)
	at org.eclipse.jetty.util.thread.QueuedThreadPool$Runner.run(QueuedThreadPool.java:1166)
	at java.base/java.lang.Thread.run(Unknown Source)
Caused by: java.security.InvalidKeyException: key length must be 32
	at java.base/sun.security.ec.ed.EdDSAPublicKeyImpl.checkLength(Unknown Source)
	at java.base/sun.security.ec.ed.EdDSAPublicKeyImpl.<init>(Unknown Source)
	at java.base/sun.security.ec.ed.EdDSAKeyFactory.generatePublicImpl(Unknown Source)
	... 118 more
~~~


# Local Environment
To pinpoint the error, a local docker environment is created, to rule out any other source. 

Only the identity hub has to be build, utilizing [the tractus-x identity hub](https://github.com/eclipse-tractusx/tractusx-identityhub).

Cloned it and then used:
~~~
./gradlew :runtimes:identityhub:dockerize
~~~
to build `identityhub:0.1.0-SNAPSHOT` which is used within the docker compose file.

> [!note] This was also tested with components based on v0.15.1 - issue is still the same. But this version is at least publicly available.
