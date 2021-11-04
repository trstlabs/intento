(window.webpackJsonp=window.webpackJsonp||[]).push([[78],{594:function(c,Z,l){"use strict";l.r(Z);var d=l(0),b=Object(d.a)({},(function(){var c=this,Z=c.$createElement,l=c._self._c||Z;return l("ContentSlotsDistributor",{attrs:{"slot-key":c.$parent.slotKey}},[l("h1",{attrs:{id:"upgrading-a-blockchain-to-use-trustless-hub-v0-18"}},[l("a",{staticClass:"header-anchor",attrs:{href:"#upgrading-a-blockchain-to-use-trustless-hub-v0-18"}},[c._v("#")]),c._v(" Upgrading a Blockchain to use Trustless Hub v0.18")]),c._v(" "),l("p",[c._v("Trustless Hub v0.18 comes with Cosmos SDK v0.44. This version of Cosmos SDK introduced changes that are not compatible with chains that were scaffolded with Trustless Hub versions lower than v0.18.")]),c._v(" "),l("p",[l("strong",[c._v("Important:")]),c._v(" After upgrading from Trustless Hub v0.17.3 to Trustless Hub v0.18, you must update the default blockchain template to use blockchains that were scaffolded with earlier versions.")]),c._v(" "),l("p",[c._v("These instructions are written for a blockchain that was scaffolded with the following command:")]),c._v(" "),l("tm-code-block",{staticClass:"codeblock",attrs:{language:"",base64:"VHJ1c3RsZXNzIEh1YiBzY2FmZm9sZCBjaGFpbiBnaXRodWIuY29tL2Nvc21vbmF1dC9tYXJzCg=="}}),c._v(" "),l("p",[c._v("If you used a different module path, replace "),l("code",[c._v("cosmonaut")]),c._v(" and "),l("code",[c._v("mars")]),c._v(" with the correct values for your blockchain.")]),c._v(" "),l("h2",{attrs:{id:"blockchain"}},[l("a",{staticClass:"header-anchor",attrs:{href:"#blockchain"}},[c._v("#")]),c._v(" Blockchain")]),c._v(" "),l("p",[c._v("For each file listed, make the required changes to the source code of the blockchain template.")]),c._v(" "),l("h3",{attrs:{id:"go-mod"}},[l("a",{staticClass:"header-anchor",attrs:{href:"#go-mod"}},[c._v("#")]),c._v(" go.mod")]),c._v(" "),l("tm-code-block",{staticClass:"codeblock",attrs:{language:"",base64:"bW9kdWxlIGdpdGh1Yi5jb20vY29zbW9uYXV0L21hcnMKCmdvIDEuMTYKCnJlcXVpcmUgKAoJZ2l0aHViLmNvbS9jb3Ntb3MvY29zbW9zLXNkayB2MC40NC4wCglnaXRodWIuY29tL2Nvc21vcy9pYmMtZ28gdjEuMi4wCglnaXRodWIuY29tL2dvZ28vcHJvdG9idWYgdjEuMy4zCglnaXRodWIuY29tL2dvb2dsZS9nby1jbXAgdjAuNS42IC8vIGluZGlyZWN0CglnaXRodWIuY29tL2dvcmlsbGEvbXV4IHYxLjguMAoJZ2l0aHViLmNvbS9ncnBjLWVjb3N5c3RlbS9ncnBjLWdhdGV3YXkgdjEuMTYuMAoJZ2l0aHViLmNvbS9zcGYxMy9jYXN0IHYxLjMuMQoJZ2l0aHViLmNvbS9zcGYxMy9jb2JyYSB2MS4xLjMKCWdpdGh1Yi5jb20vc3RyZXRjaHIvdGVzdGlmeSB2MS43LjAKCWdpdGh1Yi5jb20vdGVuZGVybWludC9zcG0gdjAuMS42CglnaXRodWIuY29tL3RlbmRlcm1pbnQvdGVuZGVybWludCB2MC4zNC4xMwoJZ2l0aHViLmNvbS90ZW5kZXJtaW50L3RtLWRiIHYwLjYuNAoJZ29vZ2xlLmdvbGFuZy5vcmcvZ2VucHJvdG8gdjAuMC4wLTIwMjEwOTAzMTYyNjQ5LWQwOGM2OGFkYmE4MwoJZ29vZ2xlLmdvbGFuZy5vcmcvZ3JwYyB2MS40MC4wCikKCnJlcGxhY2UgKAoJZ2l0aHViLmNvbS85OWRlc2lnbnMva2V5cmluZyA9Jmd0OyBnaXRodWIuY29tL2Nvc21vcy9rZXlyaW5nIHYxLjEuNy0wLjIwMjEwNjIyMTExOTEyLWVmMDBmOGFjM2Q3NgoJZ2l0aHViLmNvbS9nb2dvL3Byb3RvYnVmID0mZ3Q7IGdpdGh1Yi5jb20vcmVnZW4tbmV0d29yay9wcm90b2J1ZiB2MS4zLjMtYWxwaGEucmVnZW4uMQoJZ29vZ2xlLmdvbGFuZy5vcmcvZ3JwYyA9Jmd0OyBnb29nbGUuZ29sYW5nLm9yZy9ncnBjIHYxLjMzLjIKKQo="}}),c._v(" "),l("h3",{attrs:{id:"app-app-go"}},[l("a",{staticClass:"header-anchor",attrs:{href:"#app-app-go"}},[c._v("#")]),c._v(" app/app.go")]),c._v(" "),l("tm-code-block",{staticClass:"codeblock",attrs:{language:"go",base64:"aW1wb3J0ICgKICAvLy4uLgogIC8vIEFkZCB0aGUgZm9sbG93aW5nIHBhY2thZ2VzOgogICZxdW90O2dpdGh1Yi5jb20vY29zbW9zL2Nvc21vcy1zZGsveC9mZWVncmFudCZxdW90OwogIGZlZWdyYW50a2VlcGVyICZxdW90O2dpdGh1Yi5jb20vY29zbW9zL2Nvc21vcy1zZGsveC9mZWVncmFudC9rZWVwZXImcXVvdDsKICBmZWVncmFudG1vZHVsZSAmcXVvdDtnaXRodWIuY29tL2Nvc21vcy9jb3Ntb3Mtc2RrL3gvZmVlZ3JhbnQvbW9kdWxlJnF1b3Q7CgogICZxdW90O2dpdGh1Yi5jb20vY29zbW9zL2liYy1nby9tb2R1bGVzL2FwcHMvdHJhbnNmZXImcXVvdDsKICBpYmN0cmFuc2ZlcmtlZXBlciAmcXVvdDtnaXRodWIuY29tL2Nvc21vcy9pYmMtZ28vbW9kdWxlcy9hcHBzL3RyYW5zZmVyL2tlZXBlciZxdW90OwogIGliY3RyYW5zZmVydHlwZXMgJnF1b3Q7Z2l0aHViLmNvbS9jb3Ntb3MvaWJjLWdvL21vZHVsZXMvYXBwcy90cmFuc2Zlci90eXBlcyZxdW90OwogIGliYyAmcXVvdDtnaXRodWIuY29tL2Nvc21vcy9pYmMtZ28vbW9kdWxlcy9jb3JlJnF1b3Q7CiAgaWJjY2xpZW50ICZxdW90O2dpdGh1Yi5jb20vY29zbW9zL2liYy1nby9tb2R1bGVzL2NvcmUvMDItY2xpZW50JnF1b3Q7CiAgaWJjcG9ydHR5cGVzICZxdW90O2dpdGh1Yi5jb20vY29zbW9zL2liYy1nby9tb2R1bGVzL2NvcmUvMDUtcG9ydC90eXBlcyZxdW90OwogIGliY2hvc3QgJnF1b3Q7Z2l0aHViLmNvbS9jb3Ntb3MvaWJjLWdvL21vZHVsZXMvY29yZS8yNC1ob3N0JnF1b3Q7CiAgaWJja2VlcGVyICZxdW90O2dpdGh1Yi5jb20vY29zbW9zL2liYy1nby9tb2R1bGVzL2NvcmUva2VlcGVyJnF1b3Q7CiAgCiAgLy8gUmVtb3ZlIHRoZSBmb2xsb3dpbmcgcGFja2FnZXM6CiAgLy8gdHJhbnNmZXIgJnF1b3Q7Z2l0aHViLmNvbS9jb3Ntb3MvY29zbW9zLXNkay94L2liYy9hcHBsaWNhdGlvbnMvdHJhbnNmZXImcXVvdDsKICAvLyBpYmN0cmFuc2ZlcmtlZXBlciAmcXVvdDtnaXRodWIuY29tL2Nvc21vcy9jb3Ntb3Mtc2RrL3gvaWJjL2FwcGxpY2F0aW9ucy90cmFuc2Zlci9rZWVwZXImcXVvdDsKICAvLyBpYmN0cmFuc2ZlcnR5cGVzICZxdW90O2dpdGh1Yi5jb20vY29zbW9zL2Nvc21vcy1zZGsveC9pYmMvYXBwbGljYXRpb25zL3RyYW5zZmVyL3R5cGVzJnF1b3Q7CiAgLy8gaWJjICZxdW90O2dpdGh1Yi5jb20vY29zbW9zL2Nvc21vcy1zZGsveC9pYmMvY29yZSZxdW90OwogIC8vIGliY2NsaWVudCAmcXVvdDtnaXRodWIuY29tL2Nvc21vcy9jb3Ntb3Mtc2RrL3gvaWJjL2NvcmUvMDItY2xpZW50JnF1b3Q7CiAgLy8gcG9ydHR5cGVzICZxdW90O2dpdGh1Yi5jb20vY29zbW9zL2Nvc21vcy1zZGsveC9pYmMvY29yZS8wNS1wb3J0L3R5cGVzJnF1b3Q7CiAgLy8gaWJjaG9zdCAmcXVvdDtnaXRodWIuY29tL2Nvc21vcy9jb3Ntb3Mtc2RrL3gvaWJjL2NvcmUvMjQtaG9zdCZxdW90OwogIC8vIGliY2tlZXBlciAmcXVvdDtnaXRodWIuY29tL2Nvc21vcy9jb3Ntb3Mtc2RrL3gvaWJjL2NvcmUva2VlcGVyJnF1b3Q7CikKCnZhciAoCiAgLy8uLi4KICBNb2R1bGVCYXNpY3MgPSBtb2R1bGUuTmV3QmFzaWNNYW5hZ2VyKAogICAgLy8uLi4KICAgIHNsYXNoaW5nLkFwcE1vZHVsZUJhc2lje30sCiAgICAvLyBBZGQgZmVlZ3JhbnRtb2R1bGUuQXBwTW9kdWxlQmFzaWN7fSwKICAgIGZlZWdyYW50bW9kdWxlLkFwcE1vZHVsZUJhc2lje30sIC8vICZsdDstLQogICAgaWJjLkFwcE1vZHVsZUJhc2lje30sCiAgICAvLy4uLgogICkKICAvLy4uLgopCgp0eXBlIEFwcCBzdHJ1Y3QgewogIC8vLi4uCiAgLy8gUmVwbGFjZSBjb2RlYy5NYXJzaGFsZXIgd2l0aCBjb2RlYy5Db2RlYwogIGFwcENvZGVjICAgICAgICAgIGNvZGVjLkNvZGVjIC8vICZsdDstLQogIC8vIEFkZCBGZWVHcmFudEtlZXBlcgogIEZlZUdyYW50S2VlcGVyICAgZmVlZ3JhbnRrZWVwZXIuS2VlcGVyIC8vICZsdDstLQp9CgpmdW5jIE5ldyguLi4pIHsKICAvL2JBcHAuU2V0QXBwVmVyc2lvbih2ZXJzaW9uLlZlcnNpb24pCiAgYkFwcC5TZXRWZXJzaW9uKHZlcnNpb24uVmVyc2lvbikgLy8gJmx0Oy0tCgogIGtleXMgOj0gc2RrLk5ld0tWU3RvcmVLZXlzKAogICAgLy8uLi4KICAgIHVwZ3JhZGV0eXBlcy5TdG9yZUtleSwKICAgIC8vIEFkZCBmZWVncmFudC5TdG9yZUtleQogICAgZmVlZ3JhbnQuU3RvcmVLZXksIC8vICZsdDstLQogICAgZXZpZGVuY2V0eXBlcy5TdG9yZUtleSwKICAgIC8vLi4uCiAgKQoKICBhcHAuRmVlR3JhbnRLZWVwZXIgPSBmZWVncmFudGtlZXBlci5OZXdLZWVwZXIoYXBwQ29kZWMsIGtleXNbZmVlZ3JhbnQuU3RvcmVLZXldLCBhcHAuQWNjb3VudEtlZXBlcikgIC8vICZsdDstLQogIC8vIEFkZCBhcHAuQmFzZUFwcCBhcyB0aGUgbGFzdCBhcmd1bWVudCB0byB1cGdyYWRla2VlcGVyLk5ld0tlZXBlcgogIGFwcC5VcGdyYWRlS2VlcGVyID0gdXBncmFkZWtlZXBlci5OZXdLZWVwZXIoc2tpcFVwZ3JhZGVIZWlnaHRzLCBrZXlzW3VwZ3JhZGV0eXBlcy5TdG9yZUtleV0sIGFwcENvZGVjLCBob21lUGF0aCwgYXBwLkJhc2VBcHApCiAgCiAgYXBwLklCQ0tlZXBlciA9IGliY2tlZXBlci5OZXdLZWVwZXIoCiAgICAvLyBBZGQgYXBwLlVwZ3JhZGVLZWVwZXIKICAgIGFwcENvZGVjLCBrZXlzW2liY2hvc3QuU3RvcmVLZXldLCBhcHAuR2V0U3Vic3BhY2UoaWJjaG9zdC5Nb2R1bGVOYW1lKSwgYXBwLlN0YWtpbmdLZWVwZXIsIGFwcC5VcGdyYWRlS2VlcGVyLCBzY29wZWRJQkNLZWVwZXIsCiAgKQoKICBnb3ZSb3V0ZXIuQWRkUm91dGUoZ292dHlwZXMuUm91dGVyS2V5LCBnb3Z0eXBlcy5Qcm9wb3NhbEhhbmRsZXIpLgogICAgLy8uLi4KICAgIC8vIFJlcGxhY2UgTmV3Q2xpZW50VXBkYXRlUHJvcG9zYWxIYW5kbGVyIHdpdGggTmV3Q2xpZW50UHJvcG9zYWxIYW5kbGVyCiAgICBBZGRSb3V0ZShpYmNob3N0LlJvdXRlcktleSwgaWJjY2xpZW50Lk5ld0NsaWVudFByb3Bvc2FsSGFuZGxlcihhcHAuSUJDS2VlcGVyLkNsaWVudEtlZXBlcikpCgogIC8vIFJlcGxhY2UgcG9ydHR5cGVzIHdpdGggaWJjcG9ydHR5cGVzCiAgaWJjUm91dGVyIDo9IGliY3BvcnR0eXBlcy5OZXdSb3V0ZXIoKQoKICBhcHAubW0uU2V0T3JkZXJCZWdpbkJsb2NrZXJzKAogICAgdXBncmFkZXR5cGVzLk1vZHVsZU5hbWUsCiAgICAvLyBBZGQgY2FwYWJpbGl0eXR5cGVzLk1vZHVsZU5hbWUsCiAgICBjYXBhYmlsaXR5dHlwZXMuTW9kdWxlTmFtZSwKICAgIG1pbnR0eXBlcy5Nb2R1bGVOYW1lLAogICAgLy8uLi4KICAgIC8vIEFkZCBmZWVncmFudC5Nb2R1bGVOYW1lLAogICAgZmVlZ3JhbnQuTW9kdWxlTmFtZSwKICApCgogIC8vIEFkZCBhcHAuYXBwQ29kZWMgYXMgYW4gYXJndW1lbnQgdG8gbW9kdWxlLk5ld0NvbmZpZ3VyYXRvcjoKICBhcHAubW0uUmVnaXN0ZXJTZXJ2aWNlcyhtb2R1bGUuTmV3Q29uZmlndXJhdG9yKGFwcC5hcHBDb2RlYywgYXBwLk1zZ1NlcnZpY2VSb3V0ZXIoKSwgYXBwLkdSUENRdWVyeVJvdXRlcigpKSkKICAKICAvLyBSZXBsYWNlOgogIC8vIGFwcC5TZXRBbnRlSGFuZGxlcigKICAvLyAJYW50ZS5OZXdBbnRlSGFuZGxlcigKICAvLyAJCWFwcC5BY2NvdW50S2VlcGVyLCBhcHAuQmFua0tlZXBlciwgYW50ZS5EZWZhdWx0U2lnVmVyaWZpY2F0aW9uR2FzQ29uc3VtZXIsCiAgLy8gCQllbmNvZGluZ0NvbmZpZy5UeENvbmZpZy5TaWduTW9kZUhhbmRsZXIoKSwKICAvLyAJKSwKICAvLyApCgogIC8vIFdpdGggdGhlIGZvbGxvd2luZzoKICBhbnRlSGFuZGxlciwgZXJyIDo9IGFudGUuTmV3QW50ZUhhbmRsZXIoCiAgICBhbnRlLkhhbmRsZXJPcHRpb25zewogICAgICBBY2NvdW50S2VlcGVyOiAgIGFwcC5BY2NvdW50S2VlcGVyLAogICAgICBCYW5rS2VlcGVyOiAgICAgIGFwcC5CYW5rS2VlcGVyLAogICAgICBTaWduTW9kZUhhbmRsZXI6IGVuY29kaW5nQ29uZmlnLlR4Q29uZmlnLlNpZ25Nb2RlSGFuZGxlcigpLAogICAgICBGZWVncmFudEtlZXBlcjogIGFwcC5GZWVHcmFudEtlZXBlciwKICAgICAgU2lnR2FzQ29uc3VtZXI6ICBhbnRlLkRlZmF1bHRTaWdWZXJpZmljYXRpb25HYXNDb25zdW1lciwKICAgIH0sCiAgKQogIGlmIGVyciAhPSBuaWwgewogICAgcGFuaWMoZXJyKQogIH0KICBhcHAuU2V0QW50ZUhhbmRsZXIoYW50ZUhhbmRsZXIpCgogIC8vIFJlbW92ZSB0aGUgZm9sbG93aW5nOgogIC8vIGN0eCA6PSBhcHAuQmFzZUFwcC5OZXdVbmNhY2hlZENvbnRleHQodHJ1ZSwgdG1wcm90by5IZWFkZXJ7fSkKICAvLyBhcHAuQ2FwYWJpbGl0eUtlZXBlci5Jbml0aWFsaXplQW5kU2VhbChjdHgpCn0KCmZ1bmMgKGFwcCAqQXBwKSBJbml0Q2hhaW5lcihjdHggc2RrLkNvbnRleHQsIHJlcSBhYmNpLlJlcXVlc3RJbml0Q2hhaW4pIGFiY2kuUmVzcG9uc2VJbml0Q2hhaW4gewogIHZhciBnZW5lc2lzU3RhdGUgR2VuZXNpc1N0YXRlCiAgaWYgZXJyIDo9IHRtanNvbi5Vbm1hcnNoYWwocmVxLkFwcFN0YXRlQnl0ZXMsICZhbXA7Z2VuZXNpc1N0YXRlKTsgZXJyICE9IG5pbCB7CiAgICBwYW5pYyhlcnIpCiAgfQogIC8vIEFkZCB0aGUgZm9sbG93aW5nOgogIGFwcC5VcGdyYWRlS2VlcGVyLlNldE1vZHVsZVZlcnNpb25NYXAoY3R4LCBhcHAubW0uR2V0VmVyc2lvbk1hcCgpKQogIHJldHVybiBhcHAubW0uSW5pdEdlbmVzaXMoY3R4LCBhcHAuYXBwQ29kZWMsIGdlbmVzaXNTdGF0ZSkKfQoKLy8gUmVwbGFjZSBNYXJzaGFsZXIgd2l0aCBDb2RlYwpmdW5jIChhcHAgKkFwcCkgQXBwQ29kZWMoKSBjb2RlYy5Db2RlYyB7CiAgcmV0dXJuIGFwcC5hcHBDb2RlYwp9CgovLyBSZXBsYWNlIEJpbmFyeU1hcnNoYWxlciB3aXRoIEJpbmFyeUNvZGVjCmZ1bmMgaW5pdFBhcmFtc0tlZXBlcihhcHBDb2RlYyBjb2RlYy5CaW5hcnlDb2RlYywgbGVnYWN5QW1pbm8gKmNvZGVjLkxlZ2FjeUFtaW5vLCBrZXksIHRrZXkgc2RrLlN0b3JlS2V5KSBwYXJhbXNrZWVwZXIuS2VlcGVyIHsKICAvLy4uLgp9Cg=="}}),c._v(" "),l("h3",{attrs:{id:"app-genesis-go"}},[l("a",{staticClass:"header-anchor",attrs:{href:"#app-genesis-go"}},[c._v("#")]),c._v(" app/genesis.go")]),c._v(" "),l("tm-code-block",{staticClass:"codeblock",attrs:{language:"go",base64:"Ly8gUmVwbGFjZSBjb2RlYy5KU09OTWFyc2hhbGVyIHdpdGggY29kZWMuSlNPTkNvZGVjCmZ1bmMgTmV3RGVmYXVsdEdlbmVzaXNTdGF0ZShjZGMgY29kZWMuSlNPTkNvZGVjKSBHZW5lc2lzU3RhdGUgewogIC8vLi4uCn0K"}}),c._v(" "),l("h3",{attrs:{id:"testutil-keeper-mars-go"}},[l("a",{staticClass:"header-anchor",attrs:{href:"#testutil-keeper-mars-go"}},[c._v("#")]),c._v(" testutil/keeper/mars.go")]),c._v(" "),l("p",[c._v("Add the following code:")]),c._v(" "),l("tm-code-block",{staticClass:"codeblock",attrs:{language:"go",base64:"cGFja2FnZSBrZWVwZXIKCmltcG9ydCAoCiAgJnF1b3Q7dGVzdGluZyZxdW90OwoKICAmcXVvdDtnaXRodWIuY29tL2Nvc21vbmF1dC9tYXJzL3gvbWFycy9rZWVwZXImcXVvdDsKICAmcXVvdDtnaXRodWIuY29tL2Nvc21vbmF1dC9tYXJzL3gvbWFycy90eXBlcyZxdW90OwogICZxdW90O2dpdGh1Yi5jb20vY29zbW9zL2Nvc21vcy1zZGsvY29kZWMmcXVvdDsKICBjb2RlY3R5cGVzICZxdW90O2dpdGh1Yi5jb20vY29zbW9zL2Nvc21vcy1zZGsvY29kZWMvdHlwZXMmcXVvdDsKICAmcXVvdDtnaXRodWIuY29tL2Nvc21vcy9jb3Ntb3Mtc2RrL3N0b3JlJnF1b3Q7CiAgc3RvcmV0eXBlcyAmcXVvdDtnaXRodWIuY29tL2Nvc21vcy9jb3Ntb3Mtc2RrL3N0b3JlL3R5cGVzJnF1b3Q7CiAgc2RrICZxdW90O2dpdGh1Yi5jb20vY29zbW9zL2Nvc21vcy1zZGsvdHlwZXMmcXVvdDsKICAmcXVvdDtnaXRodWIuY29tL3N0cmV0Y2hyL3Rlc3RpZnkvcmVxdWlyZSZxdW90OwogICZxdW90O2dpdGh1Yi5jb20vdGVuZGVybWludC90ZW5kZXJtaW50L2xpYnMvbG9nJnF1b3Q7CiAgdG1wcm90byAmcXVvdDtnaXRodWIuY29tL3RlbmRlcm1pbnQvdGVuZGVybWludC9wcm90by90ZW5kZXJtaW50L3R5cGVzJnF1b3Q7CiAgdG1kYiAmcXVvdDtnaXRodWIuY29tL3RlbmRlcm1pbnQvdG0tZGImcXVvdDsKKQoKZnVuYyBNYXJzS2VlcGVyKHQgdGVzdGluZy5UQikgKCprZWVwZXIuS2VlcGVyLCBzZGsuQ29udGV4dCkgewogIHN0b3JlS2V5IDo9IHNkay5OZXdLVlN0b3JlS2V5KHR5cGVzLlN0b3JlS2V5KQogIG1lbVN0b3JlS2V5IDo9IHN0b3JldHlwZXMuTmV3TWVtb3J5U3RvcmVLZXkodHlwZXMuTWVtU3RvcmVLZXkpCgogIGRiIDo9IHRtZGIuTmV3TWVtREIoKQogIHN0YXRlU3RvcmUgOj0gc3RvcmUuTmV3Q29tbWl0TXVsdGlTdG9yZShkYikKICBzdGF0ZVN0b3JlLk1vdW50U3RvcmVXaXRoREIoc3RvcmVLZXksIHNkay5TdG9yZVR5cGVJQVZMLCBkYikKICBzdGF0ZVN0b3JlLk1vdW50U3RvcmVXaXRoREIobWVtU3RvcmVLZXksIHNkay5TdG9yZVR5cGVNZW1vcnksIG5pbCkKICByZXF1aXJlLk5vRXJyb3IodCwgc3RhdGVTdG9yZS5Mb2FkTGF0ZXN0VmVyc2lvbigpKQoKICByZWdpc3RyeSA6PSBjb2RlY3R5cGVzLk5ld0ludGVyZmFjZVJlZ2lzdHJ5KCkKICBrIDo9IGtlZXBlci5OZXdLZWVwZXIoCiAgICBjb2RlYy5OZXdQcm90b0NvZGVjKHJlZ2lzdHJ5KSwKICAgIHN0b3JlS2V5LAogICAgbWVtU3RvcmVLZXksCiAgKQoKICBjdHggOj0gc2RrLk5ld0NvbnRleHQoc3RhdGVTdG9yZSwgdG1wcm90by5IZWFkZXJ7fSwgZmFsc2UsIGxvZy5OZXdOb3BMb2dnZXIoKSkKICByZXR1cm4gaywgY3R4Cn0K"}}),c._v(" "),l("p",[c._v("If "),l("code",[c._v("mars")]),c._v(" is an IBC-enabled module, add the following code, instead:")]),c._v(" "),l("tm-code-block",{staticClass:"codeblock",attrs:{language:"go",base64:"cGFja2FnZSBrZWVwZXIKCmltcG9ydCAoCiAgJnF1b3Q7dGVzdGluZyZxdW90OwoKICAmcXVvdDtnaXRodWIuY29tL2Nvc21vbmF1dC90ZXN0L3gvbWFycy9rZWVwZXImcXVvdDsKICAmcXVvdDtnaXRodWIuY29tL2Nvc21vbmF1dC90ZXN0L3gvbWFycy90eXBlcyZxdW90OwogICZxdW90O2dpdGh1Yi5jb20vY29zbW9zL2Nvc21vcy1zZGsvY29kZWMmcXVvdDsKICBjb2RlY3R5cGVzICZxdW90O2dpdGh1Yi5jb20vY29zbW9zL2Nvc21vcy1zZGsvY29kZWMvdHlwZXMmcXVvdDsKICAmcXVvdDtnaXRodWIuY29tL2Nvc21vcy9jb3Ntb3Mtc2RrL3N0b3JlJnF1b3Q7CiAgc3RvcmV0eXBlcyAmcXVvdDtnaXRodWIuY29tL2Nvc21vcy9jb3Ntb3Mtc2RrL3N0b3JlL3R5cGVzJnF1b3Q7CiAgc2RrICZxdW90O2dpdGh1Yi5jb20vY29zbW9zL2Nvc21vcy1zZGsvdHlwZXMmcXVvdDsKICBjYXBhYmlsaXR5a2VlcGVyICZxdW90O2dpdGh1Yi5jb20vY29zbW9zL2Nvc21vcy1zZGsveC9jYXBhYmlsaXR5L2tlZXBlciZxdW90OwogIHR5cGVzcGFyYW1zICZxdW90O2dpdGh1Yi5jb20vY29zbW9zL2Nvc21vcy1zZGsveC9wYXJhbXMvdHlwZXMmcXVvdDsKICBpYmNrZWVwZXIgJnF1b3Q7Z2l0aHViLmNvbS9jb3Ntb3MvaWJjLWdvL21vZHVsZXMvY29yZS9rZWVwZXImcXVvdDsKICAmcXVvdDtnaXRodWIuY29tL3N0cmV0Y2hyL3Rlc3RpZnkvcmVxdWlyZSZxdW90OwogICZxdW90O2dpdGh1Yi5jb20vdGVuZGVybWludC90ZW5kZXJtaW50L2xpYnMvbG9nJnF1b3Q7CiAgdG1wcm90byAmcXVvdDtnaXRodWIuY29tL3RlbmRlcm1pbnQvdGVuZGVybWludC9wcm90by90ZW5kZXJtaW50L3R5cGVzJnF1b3Q7CiAgdG1kYiAmcXVvdDtnaXRodWIuY29tL3RlbmRlcm1pbnQvdG0tZGImcXVvdDsKKQoKZnVuYyBNYXJzS2VlcGVyKHQgdGVzdGluZy5UQikgKCprZWVwZXIuS2VlcGVyLCBzZGsuQ29udGV4dCkgewogIGxvZ2dlciA6PSBsb2cuTmV3Tm9wTG9nZ2VyKCkKCiAgc3RvcmVLZXkgOj0gc2RrLk5ld0tWU3RvcmVLZXkodHlwZXMuU3RvcmVLZXkpCiAgbWVtU3RvcmVLZXkgOj0gc3RvcmV0eXBlcy5OZXdNZW1vcnlTdG9yZUtleSh0eXBlcy5NZW1TdG9yZUtleSkKCiAgZGIgOj0gdG1kYi5OZXdNZW1EQigpCiAgc3RhdGVTdG9yZSA6PSBzdG9yZS5OZXdDb21taXRNdWx0aVN0b3JlKGRiKQogIHN0YXRlU3RvcmUuTW91bnRTdG9yZVdpdGhEQihzdG9yZUtleSwgc2RrLlN0b3JlVHlwZUlBVkwsIGRiKQogIHN0YXRlU3RvcmUuTW91bnRTdG9yZVdpdGhEQihtZW1TdG9yZUtleSwgc2RrLlN0b3JlVHlwZU1lbW9yeSwgbmlsKQogIHJlcXVpcmUuTm9FcnJvcih0LCBzdGF0ZVN0b3JlLkxvYWRMYXRlc3RWZXJzaW9uKCkpCgogIHJlZ2lzdHJ5IDo9IGNvZGVjdHlwZXMuTmV3SW50ZXJmYWNlUmVnaXN0cnkoKQogIGFwcENvZGVjIDo9IGNvZGVjLk5ld1Byb3RvQ29kZWMocmVnaXN0cnkpCiAgY2FwYWJpbGl0eUtlZXBlciA6PSBjYXBhYmlsaXR5a2VlcGVyLk5ld0tlZXBlcihhcHBDb2RlYywgc3RvcmVLZXksIG1lbVN0b3JlS2V5KQoKICBhbWlubyA6PSBjb2RlYy5OZXdMZWdhY3lBbWlubygpCiAgc3MgOj0gdHlwZXNwYXJhbXMuTmV3U3Vic3BhY2UoYXBwQ29kZWMsCiAgICBhbWlubywKICAgIHN0b3JlS2V5LAogICAgbWVtU3RvcmVLZXksCiAgICAmcXVvdDtNYXJzU3ViU3BhY2UmcXVvdDssCiAgKQogIElCQ0tlZXBlciA6PSBpYmNrZWVwZXIuTmV3S2VlcGVyKAogICAgYXBwQ29kZWMsCiAgICBzdG9yZUtleSwKICAgIHNzLAogICAgbmlsLAogICAgbmlsLAogICAgY2FwYWJpbGl0eUtlZXBlci5TY29wZVRvTW9kdWxlKCZxdW90O01hcnNJQkNLZWVwZXImcXVvdDspLAogICkKCiAgayA6PSBrZWVwZXIuTmV3S2VlcGVyKAogICAgY29kZWMuTmV3UHJvdG9Db2RlYyhyZWdpc3RyeSksCiAgICBzdG9yZUtleSwKICAgIG1lbVN0b3JlS2V5LAogICAgSUJDS2VlcGVyLkNoYW5uZWxLZWVwZXIsCiAgICAmYW1wO0lCQ0tlZXBlci5Qb3J0S2VlcGVyLAogICAgY2FwYWJpbGl0eUtlZXBlci5TY29wZVRvTW9kdWxlKCZxdW90O01hcnNTY29wZWRLZWVwZXImcXVvdDspLAogICkKCiAgY3R4IDo9IHNkay5OZXdDb250ZXh0KHN0YXRlU3RvcmUsIHRtcHJvdG8uSGVhZGVye30sIGZhbHNlLCBsb2dnZXIpCiAgcmV0dXJuIGssIGN0eAp9Cg=="}}),c._v(" "),l("h3",{attrs:{id:"testutil-network-network-go"}},[l("a",{staticClass:"header-anchor",attrs:{href:"#testutil-network-network-go"}},[c._v("#")]),c._v(" testutil/network/network.go")]),c._v(" "),l("tm-code-block",{staticClass:"codeblock",attrs:{language:"go",base64:"ZnVuYyBEZWZhdWx0Q29uZmlnKCkgbmV0d29yay5Db25maWcgewogIC8vLi4uCiAgcmV0dXJuIG5ldHdvcmsuQ29uZmlnewogICAgLy8uLi4KICAgIC8vIEFkZCBzZGsuRGVmYXVsdFBvd2VyUmVkdWN0aW9uCiAgICBBY2NvdW50VG9rZW5zOiAgIHNkay5Ub2tlbnNGcm9tQ29uc2Vuc3VzUG93ZXIoMTAwMCwgc2RrLkRlZmF1bHRQb3dlclJlZHVjdGlvbiksCiAgICBTdGFraW5nVG9rZW5zOiAgIHNkay5Ub2tlbnNGcm9tQ29uc2Vuc3VzUG93ZXIoNTAwLCBzZGsuRGVmYXVsdFBvd2VyUmVkdWN0aW9uKSwKICAgIEJvbmRlZFRva2VuczogICAgc2RrLlRva2Vuc0Zyb21Db25zZW5zdXNQb3dlcigxMDAsIHNkay5EZWZhdWx0UG93ZXJSZWR1Y3Rpb24pLAogICAgLy8uLi4KICB9Cn0K"}}),c._v(" "),l("h3",{attrs:{id:"testutil-sample-sample-go"}},[l("a",{staticClass:"header-anchor",attrs:{href:"#testutil-sample-sample-go"}},[c._v("#")]),c._v(" testutil/sample/sample.go")]),c._v(" "),l("p",[c._v("Add the following code:")]),c._v(" "),l("tm-code-block",{staticClass:"codeblock",attrs:{language:"go",base64:"cGFja2FnZSBzYW1wbGUKCmltcG9ydCAoCiAgJnF1b3Q7Z2l0aHViLmNvbS9jb3Ntb3MvY29zbW9zLXNkay9jcnlwdG8va2V5cy9lZDI1NTE5JnF1b3Q7CiAgc2RrICZxdW90O2dpdGh1Yi5jb20vY29zbW9zL2Nvc21vcy1zZGsvdHlwZXMmcXVvdDsKKQoKLy8gQWNjQWRkcmVzcyByZXR1cm5zIGEgc2FtcGxlIGFjY291bnQgYWRkcmVzcwpmdW5jIEFjY0FkZHJlc3MoKSBzdHJpbmcgewogIHBrIDo9IGVkMjU1MTkuR2VuUHJpdktleSgpLlB1YktleSgpCiAgYWRkciA6PSBway5BZGRyZXNzKCkKICByZXR1cm4gc2RrLkFjY0FkZHJlc3MoYWRkcikuU3RyaW5nKCkKfQo="}}),c._v(" "),l("h3",{attrs:{id:"bandchain-support"}},[l("a",{staticClass:"header-anchor",attrs:{href:"#bandchain-support"}},[c._v("#")]),c._v(" BandChain Support")]),c._v(" "),l("p",[c._v("If your module includes integration with BandChain, added manually or scaffolded with "),l("code",[c._v("Trustless Hub scaffold band")]),c._v(", upgrade the "),l("code",[c._v("github.com/bandprotocol/bandchain-packet")]),c._v(" package to "),l("code",[c._v("v0.0.2")]),c._v(" in "),l("code",[c._v("go.mod")]),c._v(".")]),c._v(" "),l("h2",{attrs:{id:"module"}},[l("a",{staticClass:"header-anchor",attrs:{href:"#module"}},[c._v("#")]),c._v(" Module")]),c._v(" "),l("h3",{attrs:{id:"x-mars-keeper-keeper-go"}},[l("a",{staticClass:"header-anchor",attrs:{href:"#x-mars-keeper-keeper-go"}},[c._v("#")]),c._v(" x/mars/keeper/keeper.go")]),c._v(" "),l("tm-code-block",{staticClass:"codeblock",attrs:{language:"go",base64:"dHlwZSAoCiAgS2VlcGVyIHN0cnVjdCB7CiAgICAvLyBSZXBsYWNlIE1hcnNoYWxlciB3aXRoIEJpbmFyeUNvZGVjCiAgICBjZGMgICAgICBjb2RlYy5CaW5hcnlDb2RlYwogICAgLy8uLi4KICB9CikKCmZ1bmMgTmV3S2VlcGVyKAogIC8vIFJlcGxhY2UgTWFyc2hhbGVyIHdpdGggQmluYXJ5Q29kZWMKICBjZGMgY29kZWMuQmluYXJ5Q29kZWMsCiAgLy8uLi4KKSAqS2VlcGVyIHsKICAvLyAuLi4KfQo="}}),c._v(" "),l("h3",{attrs:{id:"x-mars-keeper-msg-server-test-go"}},[l("a",{staticClass:"header-anchor",attrs:{href:"#x-mars-keeper-msg-server-test-go"}},[c._v("#")]),c._v(" x/mars/keeper/msg_server_test.go")]),c._v(" "),l("tm-code-block",{staticClass:"codeblock",attrs:{language:"go",base64:"cGFja2FnZSBrZWVwZXJfdGVzdAoKaW1wb3J0ICgKICAvLy4uLgogIC8vIEFkZCB0aGUgZm9sbG93aW5nOgogIGtlZXBlcnRlc3QgJnF1b3Q7Z2l0aHViLmNvbS9jb3Ntb25hdXQvbWFycy90ZXN0dXRpbC9rZWVwZXImcXVvdDsKICAmcXVvdDtnaXRodWIuY29tL2Nvc21vbmF1dC9tYXJzL3gvbWFycy9rZWVwZXImcXVvdDsKKQoKZnVuYyBzZXR1cE1zZ1NlcnZlcih0IHRlc3RpbmcuVEIpICh0eXBlcy5Nc2dTZXJ2ZXIsIGNvbnRleHQuQ29udGV4dCkgewogIC8vIFJlcGxhY2UKICAvLyBrZWVwZXIsIGN0eCA6PSBzZXR1cEtlZXBlcih0KQogIC8vIHJldHVybiBOZXdNc2dTZXJ2ZXJJbXBsKCprZWVwZXIpLCBzZGsuV3JhcFNES0NvbnRleHQoY3R4KQoKICAvLyBXaXRoIHRoZSBmb2xsb3dpbmc6CiAgaywgY3R4IDo9IGtlZXBlcnRlc3QuTWFyc0tlZXBlcih0KQogIHJldHVybiBrZWVwZXIuTmV3TXNnU2VydmVySW1wbCgqayksIHNkay5XcmFwU0RLQ29udGV4dChjdHgpCn0K"}}),c._v(" "),l("h3",{attrs:{id:"x-mars-module-go"}},[l("a",{staticClass:"header-anchor",attrs:{href:"#x-mars-module-go"}},[c._v("#")]),c._v(" x/mars/module.go")]),c._v(" "),l("tm-code-block",{staticClass:"codeblock",attrs:{language:"go",base64:"dHlwZSBBcHBNb2R1bGVCYXNpYyBzdHJ1Y3QgewogIC8vIFJlcGxhY2UgTWFyc2hhbGVyIHdpdGggQmluYXJ5Q29kZWMKICBjZGMgY29kZWMuQmluYXJ5Q29kZWMKfQoKLy8gUmVwbGFjZSBNYXJzaGFsZXIgd2l0aCBCaW5hcnlDb2RlYwpmdW5jIE5ld0FwcE1vZHVsZUJhc2ljKGNkYyBjb2RlYy5CaW5hcnlDb2RlYykgQXBwTW9kdWxlQmFzaWMgewogIHJldHVybiBBcHBNb2R1bGVCYXNpY3tjZGM6IGNkY30KfQoKLy8gUmVwbGFjZSBKU09OTWFyc2hhbGVyIHdpdGggSlNPTkNvZGVjCmZ1bmMgKEFwcE1vZHVsZUJhc2ljKSBEZWZhdWx0R2VuZXNpcyhjZGMgY29kZWMuSlNPTkNvZGVjKSBqc29uLlJhd01lc3NhZ2UgewogIHJldHVybiBjZGMuTXVzdE1hcnNoYWxKU09OKHR5cGVzLkRlZmF1bHRHZW5lc2lzKCkpCn0KCi8vIFJlcGxhY2UgSlNPTk1hcnNoYWxlciB3aXRoIEpTT05Db2RlYwpmdW5jIChBcHBNb2R1bGVCYXNpYykgVmFsaWRhdGVHZW5lc2lzKGNkYyBjb2RlYy5KU09OQ29kZWMsIGNvbmZpZyBjbGllbnQuVHhFbmNvZGluZ0NvbmZpZywgYnoganNvbi5SYXdNZXNzYWdlKSBlcnJvciB7CiAgLy8uLi4KfQoKLy8gUmVwbGFjZSBjb2RlYy5NYXJzaGFsbGVyIHdpdGggY29kZWMuQ29kZWMKZnVuYyBOZXdBcHBNb2R1bGUoY2RjIGNvZGVjLkNvZGVjLCBrZWVwZXIga2VlcGVyLktlZXBlcikgQXBwTW9kdWxlIHsKICAvLy4uLgp9CgovLyBSZXBsYWNlIEpTT05NYXJzaGFsZXIgd2l0aCBKU09OQ29kZWMKZnVuYyAoYW0gQXBwTW9kdWxlKSBJbml0R2VuZXNpcyhjdHggc2RrLkNvbnRleHQsIGNkYyBjb2RlYy5KU09OQ29kZWMsIGdzIGpzb24uUmF3TWVzc2FnZSkgW11hYmNpLlZhbGlkYXRvclVwZGF0ZSB7CiAgLy8uLi4KfQoKLy8gUmVwbGFjZSBKU09OTWFyc2hhbGVyIHdpdGggSlNPTkNvZGVjCmZ1bmMgKGFtIEFwcE1vZHVsZSkgRXhwb3J0R2VuZXNpcyhjdHggc2RrLkNvbnRleHQsIGNkYyBjb2RlYy5KU09OQ29kZWMpIGpzb24uUmF3TWVzc2FnZSB7CiAgLy8uLi4KfQoKLy8gQWRkIHRoZSBmb2xsb3dpbmcKZnVuYyAoQXBwTW9kdWxlKSBDb25zZW5zdXNWZXJzaW9uKCkgdWludDY0IHsgcmV0dXJuIDIgfQo="}})],1)}),[],!1,null,null,null);Z.default=b.exports}}]);