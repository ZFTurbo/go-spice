for /l  %%x IN (0,1, 1133) DO (
go run ./cmd/pgsim -f ../SPICE-GNN/assets/fold3_big/ibmpg%%x/ibmpg%%x.spice -o ../SPICE-GNN/assets/fold3_big -ff .csv -p 1e-8
)